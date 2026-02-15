package v1

import (
	proto "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/data"
	"Advanced_Shop/app/xshop/api/internal/domain/request/order"
	"Advanced_Shop/app/xshop/api/internal/service"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type orderController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewGoodsController(srv service.ServiceFactory, trans ut.Translator) *orderController {
	return &orderController{
		srv:   srv,
		trans: trans,
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
	_, err = oc.srv.Order().SubmitOrder(ctx, &proto.OrderRequest{
		UserId:  userID,
		OrderSn: RandomSns(userID),
		Address: cr.Address,
		Name:    cr.Name,
		Mobile:  cr.Mobile,
		Post:    cr.Post,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.Config.Alipay.NotifyUrl
	p.ReturnURL = global.Config.Alipay.ReturnUrl
	p.Subject = "下次一定_" + orderModel.OrderSn
	p.OutTradeNo = orderModel.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(orderModel.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	p.TimeoutExpress = "30m" //3分钟 链接失效

	result, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "生成支付宝url失败")
		return
	}
	response := order.OrderCreateResponse{
		Id:        orderModel.Id,
		AlipayUrl: result.String(),
	}
	core.OkWithData(c, response)

}

func (oc orderController) OrderDetailView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr order_srv.OrderIdRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	userId := claims.UserID
	if claims.Role == enum.AdminRole {
		userId = 0
	}
	orderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	result, err := orderClient.OrderDetail(context.WithValue(context.Background(), "ginContext", c), &proto.OrderRequest{
		UserId: userId,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	response := order_srv.OrderDetailResponse{
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
	var goodsInfo []order_srv.GoodInfo
	for _, good := range result.Goods {
		info := order_srv.GoodInfo{
			Id:    good.Id,
			Name:  good.GoodsName,
			Image: good.GoodsImage,
			Price: good.GoodsPrice,
			Nums:  good.Nums,
		}
		goodsInfo = append(goodsInfo, info)
	}
	response.GoodInfo = goodsInfo

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.Config.Alipay.NotifyUrl
	p.ReturnURL = global.Config.Alipay.ReturnUrl
	p.Subject = "下次一定_" + response.OrderSn
	p.OutTradeNo = response.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(response.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	rep, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "生成支付宝url失败")
		return
	}
	response.AlipayUrl = rep.String()

	res.OkWithData(c, response)

}

func (oc orderController) AlipayCallBackView(c *gin.Context) {

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	notification, err := client.GetTradeNotification(c.Request)
	if err != nil || notification == nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "")
		return

	}
	OrderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	// 查看订单是否存在 是否超时
	info, err := OrderClient.OrderDetailByOrderSn(ctx, &proto.AlipayOrderSnRequest{
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
	_, err = OrderClient.UpdateOrderStatus(ctx, &proto.OrderStatus{
		OrderSn: notification.OutTradeNo,
		Status:  string(notification.TradeStatus),
	})
	if err != nil {
		c.String(200, "fail")
		return
	}
	c.String(http.StatusOK, "success")

}
