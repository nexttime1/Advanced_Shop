package es

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data_search/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/options"
	"Advanced_Shop/gnova/code"
	"Advanced_Shop/pkg/db"
	"Advanced_Shop/pkg/errors"
	zlog "Advanced_Shop/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/olivere/elastic/v7"
	"strconv"
	"sync"
	"time"
)

var (
	searchFactory v1.SearchFactory
	once          sync.Once
)

type dataSearch struct {
	esClient  *elastic.Client
	consumer  rocketmq.PushConsumer
	mqOpts    *options.RocketMQOptions
	isRunning bool // 标记消费者是否运行中
	runLock   sync.RWMutex
}

func (ds *dataSearch) Goods() v1.GoodsStore {
	return newGoods(ds)
}

func GetSearchFactoryOr(opts *options.EsOptions, mqOpts *options.RocketMQOptions) (v1.SearchFactory, error) {
	if opts == nil && searchFactory == nil {
		return nil, errors.New("failed to get es client")
	}

	once.Do(func() {
		esOpt := db.EsOptions{
			Host: opts.Host,
			Port: opts.Port,
		}
		esClient, err := db.NewEsClient(&esOpt)
		if err != nil {
			return
		}

		searchFactory = &dataSearch{esClient: esClient, mqOpts: mqOpts}
	})
	if searchFactory == nil {
		return nil, errors.New("failed to get es client")
	}
	return searchFactory, nil
}

// Listen 启动商品MQ消费者（优化版）
func (c *dataSearch) Listen(ctx context.Context) error {
	c.runLock.Lock()
	if c.isRunning {
		c.runLock.Unlock()
		return errors.WithCode(code.ErrInvalidOperation, "商品MQ消费者已在运行中")
	}
	c.isRunning = true
	c.runLock.Unlock()

	// 创建PushConsumer
	consumerIns, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{c.mqOpts.Addr()}),
		consumer.WithGroupName(c.mqOpts.ConsumerGroupName),
		// 、重试次数、最大并发等
		consumer.WithMaxReconsumeTimes(3),                             // 最大重试次数
		consumer.WithConsumeMessageBatchMaxSize(1),                    // 批量消费大小（单条消费更安全）
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset), // 从最后偏移量开始消费
	)
	if err != nil {
		errMsg := fmt.Sprintf("创建商品MQ消费者失败: %v", err)
		zlog.Error(errMsg)
		return errors.WithCode(code.ErrConnectMQ, errMsg)
	}

	// 订阅Topic
	err = consumerIns.Subscribe(
		c.mqOpts.ConsumerTopic,
		consumer.MessageSelector{},
		c.SyncGoodsToES, // 消费函数：同步商品数据到ES
	)
	if err != nil {
		errMsg := fmt.Sprintf("订阅商品MQ Topic失败: %v", err)
		zlog.Error(errMsg)
		_ = consumerIns.Shutdown() // 失败时关闭已创建的消费者
		return errors.WithCode(code.ErrConnectMQ, errMsg)
	}

	// 3. 启动消费者
	if err = consumerIns.Start(); err != nil {
		errMsg := fmt.Sprintf("启动商品MQ消费者失败: %v", err)
		zlog.Error(errMsg)
		_ = consumerIns.Shutdown()
		return errors.WithCode(code.ErrConnectMQ, errMsg)
	}

	c.consumer = consumerIns
	zlog.Infof("商品MQ消费者启动成功,topic: %v, group: %v ", c.mqOpts.ConsumerTopic, c.mqOpts.ConsumerGroupName)

	go func() {
		<-ctx.Done()
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

// Close 优雅关闭消费者
func (c *dataSearch) Close() error {
	c.runLock.Lock()
	defer c.runLock.Unlock()

	if !c.isRunning || c.consumer == nil {
		zlog.Warn("商品MQ消费者未运行，无需关闭")
		return nil
	}

	// 标记为停止状态
	c.isRunning = false

	// 优雅关闭：设置超时时间，避免阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	// 异步关闭，避免阻塞主线程
	done := make(chan error, 1)
	go func() {
		done <- c.consumer.Shutdown()
		close(done)
	}()

	select {
	case err = <-done:
		if err != nil {
			errMsg := fmt.Sprintf("关闭商品MQ消费者失败: %v", err)
			zlog.Error(errMsg)
			return errors.WithCode(code.ErrConnectMQ, errMsg)
		}
	case <-ctx.Done():
		errMsg := "关闭商品MQ消费者超时（10s），强制退出"
		zlog.Error(errMsg)
		return errors.WithCode(code.ErrConnectMQ, errMsg)
	}

	c.consumer = nil
	zlog.Info("商品MQ消费者已优雅关闭")
	return nil
}

// SyncGoodsToES 消费MQ消息，同步商品数据到ES（核心消费逻辑）
func (c *dataSearch) SyncGoodsToES(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	zlog.Info("收到商品MQ消息，开始同步到ES")

	// 遍历处理每条消息
	for _, msg := range msgs {
		zlog.Debugf("解析商品MQ消息, msgID: %v, body: %v", msg.MsgId, msg.Body)

		// 解析MQ消息体
		var msgBody struct {
			EventType string                 `json:"event_type"` // INSERT/UPDATE/DELETE
			Table     string                 `json:"table"`
			Goods     map[string]interface{} `json:"goods"` // 商品字段（id/name/category_id等）
		}
		err := json.Unmarshal(msg.Body, &msgBody)
		if err != nil {
			errMsg := fmt.Sprintf("解析商品MQ消息失败, msgID=%s, err=%v", msg.MsgId, err)
			zlog.Error(errMsg)
			// 解析失败重试
			return consumer.ConsumeRetryLater, errors.WithCode(code.ErrDecodingJSON, errMsg)
		}

		// 过滤非商品表消息（防御性校验）
		if msgBody.Table != "goods" {
			zlog.Warnf("非商品表消息，忽略 table: %v", msgBody.Table)
			continue
		}

		// 转换为ES的GoodsSearchDO
		goodsSearchDO, err := convertToGoodsSearchDO(msgBody.Goods)
		if err != nil {
			errMsg := fmt.Sprintf("转换商品数据失败, msgID=%s, err=%v", msg.MsgId, err)
			zlog.Error(errMsg)
			return consumer.ConsumeRetryLater, errors.WithCode(code.ErrValidation, errMsg)
		}

		// 根据事件类型同步到ES
		goodsStore := c.Goods() // 从SearchFactory获取ES操作实例
		switch msgBody.EventType {
		case "INSERT", "UPDATE":
			// 新增/更新：写入ES
			err = goodsStore.Create(ctx, goodsSearchDO)
			if err != nil {
				errMsg := fmt.Sprintf("写入ES失败, goodsID=%d, err=%v", goodsSearchDO.ID, err)
				zlog.Error(errMsg)
				return consumer.ConsumeRetryLater, errors.WithCode(code.ErrDatabase, errMsg)
			}
			zlog.Infof("商品数据同步到ES成功 goodsID: %v", goodsSearchDO.ID)

		case "DELETE":
			// 删除：从ES删除
			err = goodsStore.Delete(ctx, uint64(goodsSearchDO.ID))
			if err != nil {
				errMsg := fmt.Sprintf("从ES删除商品失败, goodsID=%d, err=%v", goodsSearchDO.ID, err)
				zlog.Error(errMsg)
				return consumer.ConsumeRetryLater, errors.WithCode(code.ErrDatabase, errMsg)
			}
			zlog.Infof("商品数据从ES删除成功 ID: %v", goodsSearchDO.ID)

		default:
			zlog.Warnf("不支持的事件类型，忽略, eventType : %v", msgBody.EventType)
			continue
		}
	}

	// 所有消息处理成功
	return consumer.ConsumeSuccess, nil
}

// convertToGoodsSearchDO 将MQ消息中的商品map转换为ES的GoodsSearchDO
func convertToGoodsSearchDO(goodsMap map[string]interface{}) (*do.GoodsSearchDO, error) {
	goodsDO := &do.GoodsSearchDO{}

	// 解析ID（必传字段）
	idStr, ok := goodsMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("商品ID字段缺失或类型错误")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("商品ID转换失败: %v", err)
	}
	goodsDO.ID = int32(id)

	// 解析分类ID
	if categoryIDStr, ok := goodsMap["category_id"].(string); ok {
		categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 64)
		goodsDO.CategoryID = int32(categoryID)
	}

	// 解析品牌ID
	if brandsIDStr, ok := goodsMap["brands_id"].(string); ok {
		brandsID, _ := strconv.ParseUint(brandsIDStr, 10, 64)
		goodsDO.BrandsID = int32(brandsID)
	}

	// 解析商品名称
	if name, ok := goodsMap["name"].(string); ok {
		goodsDO.Name = name
	}

	// 解析点击数
	if clickNumStr, ok := goodsMap["click_num"].(string); ok {
		clickNum, _ := strconv.ParseInt(clickNumStr, 10, 64)
		goodsDO.ClickNum = int32(clickNum)
	}

	// 解析收藏数
	if favNumStr, ok := goodsMap["fav_num"].(string); ok {
		favNum, _ := strconv.ParseInt(favNumStr, 10, 64)
		goodsDO.FavNum = int32(favNum)
	}

	// 解析市场价
	if marketPriceStr, ok := goodsMap["market_price"].(string); ok {
		marketPrice, _ := strconv.ParseFloat(marketPriceStr, 64)
		goodsDO.MarketPrice = float32(marketPrice)
	}

	// 解析店铺价
	if shopPriceStr, ok := goodsMap["shop_price"].(string); ok {
		shopPrice, _ := strconv.ParseFloat(shopPriceStr, 64)
		goodsDO.ShopPrice = float32(shopPrice)
	}

	// 解析商品简介
	if goodsBrief, ok := goodsMap["goods_brief"].(string); ok {
		goodsDO.GoodsBrief = goodsBrief
	}

	// 解析是否上架
	if onSaleStr, ok := goodsMap["on_sale"].(string); ok {
		onSale, _ := strconv.ParseBool(onSaleStr)
		goodsDO.OnSale = onSale
	}

	// 解析是否包邮
	if shipFreeStr, ok := goodsMap["ship_free"].(string); ok {
		shipFree, _ := strconv.ParseBool(shipFreeStr)
		goodsDO.ShipFree = shipFree
	}

	// 解析是否新品
	if isNewStr, ok := goodsMap["is_new"].(string); ok {
		isNew, _ := strconv.ParseBool(isNewStr)
		goodsDO.IsNew = isNew
	}

	// 解析是否热销
	if isHotStr, ok := goodsMap["is_hot"].(string); ok {
		isHot, _ := strconv.ParseBool(isHotStr)
		goodsDO.IsHot = isHot
	}

	return goodsDO, nil
}
