package v1

import (
	proto "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/pkg/common"
	"Advanced_Shop/app/pkg/options"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/order"
	"Advanced_Shop/app/xshop/api/internal/service"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"fmt"
	"github.com/smartwalle/alipay/v3"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"

	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type orderController struct {
	trans   ut.Translator
	srv     service.ServiceFactory
	options *options.AliyunOptions
}

func NewOrderController(srv service.ServiceFactory, trans ut.Translator, options *options.AliyunOptions) *orderController {
	return &orderController{
		srv:     srv,
		trans:   trans,
		options: options,
	}
}

func (oc orderController) OrderListView(c *gin.Context) {
	log.Info("order list function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	var cr common.PageInfo
	err = c.ShouldBindQuery(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, oc.trans)
		return
	}
	ctx := c.Request.Context()

	list, err := oc.srv.Order().OrderList(ctx, &proto.OrderFilterRequest{
		UserId:      userID,
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})

	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	var response []order.OrderListResponse
	for _, model := range list.Data {
		result := order.OrderListResponse{
			Id:      model.Id,
			UserId:  model.UserId,
			OrderSn: model.OrderSn,
			PayType: model.PayType,
			Status:  model.Status,
			Post:    model.Post,
			Total:   model.Total,
			Address: model.Address,
			Name:    model.Name,
			Mobile:  model.Mobile,
		}
		response = append(response, result)

	}

	core.OkWithList(c, response, list.Total)

}

func RandomSns(userID int32) string {
	now := time.Now()
	rand.Seed(now.UnixNano()) //毫秒级
	id := rand.Intn(90) + 10  // 两位随机数
	OrderSns := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userID, id)

	return OrderSns
}
func (oc orderController) OrderCreateView(c *gin.Context) {
	log.Info("order create function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	var cr order.OrderCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, oc.trans)
		return
	}
	ctx := c.Request.Context()
	orderSn := RandomSns(userID)
	total, err := oc.srv.Order().SubmitOrder(ctx, &proto.OrderRequest{
		UserId:  userID,
		OrderSn: orderSn,
		Address: cr.Address,
		Name:    cr.Name,
		Mobile:  cr.Mobile,
		Post:    cr.Post,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	client, err := alipay.New(oc.options.AlipayAppId, oc.options.AlipayPrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(oc.options.AlipayPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = oc.options.AlipayNotifyUrl
	p.ReturnURL = oc.options.AlipayReturnUrl
	p.Subject = oc.options.AlipaySubject + orderSn
	p.OutTradeNo = orderSn
	p.TotalAmount = strconv.FormatFloat(float64(total.PriceSum), 'f', 2, 64)
	p.ProductCode = oc.options.AlipayProductCode
	p.TimeoutExpress = oc.options.AlipayTimeoutExpress // 默认3分钟 链接失效

	result, err := client.TradePagePay(p)
	if err != nil {
		log.Errorf("生成支付宝url失败")
		core.WriteErrResponse(c, err, nil)
		return
	}
	response := order.OrderCreateResponse{
		OrderSn:   orderSn,
		AlipayUrl: result.String(),
	}

	core.OkWithData(c, response)

}

func (oc orderController) OrderDetailView(c *gin.Context) {
	log.Info("order detail function called ...")
	userID, role, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	var cr order.OrderIdRequest
	err = c.ShouldBindQuery(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, oc.trans)
		return
	}

	if role == 1 {
		userID = 0
	}
	ctx := c.Request.Context()
	result, err := oc.srv.Order().OrderDetail(ctx, &proto.OrderRequest{
		UserId: userID,
	})
	if err != nil {
		return
	}

	response := order.OrderDetailResponse{
		Id:      result.OrderInfo.Id,
		UserId:  result.OrderInfo.UserId,
		OrderSn: result.OrderInfo.OrderSn,
		PayType: result.OrderInfo.PayType,
		Status:  result.OrderInfo.Status,
		Post:    result.OrderInfo.Post,
		Total:   result.OrderInfo.Total,
		Address: result.OrderInfo.Address,
		Name:    result.OrderInfo.Name,
		Mobile:  result.OrderInfo.Mobile,
	}
	var goodsInfo []order.GoodInfo
	for _, good := range result.Goods {
		info := order.GoodInfo{
			Id:    good.Id,
			Name:  good.GoodsName,
			Image: good.GoodsImage,
			Price: good.GoodsPrice,
			Nums:  good.Nums,
		}
		goodsInfo = append(goodsInfo, info)
	}
	response.GoodInfo = goodsInfo
	client, err := alipay.New(oc.options.AlipayAppId, oc.options.AlipayPrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(oc.options.AlipayPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = oc.options.AlipayNotifyUrl
	p.ReturnURL = oc.options.AlipayReturnUrl
	p.Subject = oc.options.AlipaySubject + result.OrderInfo.OrderSn
	p.OutTradeNo = result.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(result.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = oc.options.AlipayProductCode
	p.TimeoutExpress = oc.options.AlipayTimeoutExpress // 默认3分钟 链接失效

	alipayRes, err := client.TradePagePay(p)

	if err != nil {
		log.Errorf("生成支付宝url失败")
		core.WriteErrResponse(c, err, nil)
		return
	}
	response.AlipayUrl = alipayRes.String()

	core.OkWithData(c, response)

}

func (oc orderController) AlipayCallBackView(c *gin.Context) {

	client, err := alipay.New(oc.options.AlipayAppId, oc.options.AlipayPrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(oc.options.AlipayPublicKey)
	if err != nil {
		panic(err)
	}

	notification, err := client.GetTradeNotification(c.Request)
	if err != nil || notification == nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	// 查看订单是否存在 是否超时
	ctx := c.Request.Context()
	info, err := oc.srv.Order().OrderDetailByOrderSn(ctx, &proto.AlipayOrderSnRequest{
		OrderSn: notification.OutTradeNo,
	})
	if err != nil {
		zap.S().Error(err)
		c.String(200, "fail")
		return
	}
	if info.OrderInfo.Status != "" {
		c.String(200, "fail")
	}
	_, err = oc.srv.Order().UpdateOrderStatus(ctx, &proto.OrderStatus{
		OrderSn: notification.OutTradeNo,
		Status:  string(notification.TradeStatus),
	})
	if err != nil {
		c.String(200, "fail")
		return
	}
	c.String(http.StatusOK, "success")

}
