package v1

import (
	proto "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/order"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (oc orderController) CartListView(c *gin.Context) {
	log.Info("Cart list function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	ctx := c.Request.Context()
	response, total, err := oc.srv.Order().CartItemList(ctx, &proto.UserInfo{
		Id: userID,
	})

	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	core.OkWithList(c, response, total)

}

func (oc orderController) DeleteCartItemView(c *gin.Context) {
	log.Info("DeleteCartItemView function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr order.CartIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, nil)
		return
	}
	ctx := c.Request.Context()
	_, err = oc.srv.Order().DeleteCartItem(ctx, &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	core.OkWithMessage(c, "删除成功")

}

func (oc orderController) AddItemView(c *gin.Context) {
	log.Info("AddItemView function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr order.CartAddRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, nil)
		return
	}
	ctx := c.Request.Context()
	oc.srv.Order().CreateCartItem(ctx, &proto.CartItemRequest{})
	res.OkWithData(c, response)

}

func (oc orderController) UpdatePatchView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr cart_srv.CartIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
	}
	var update cart_srv.CartUpdateRequest
	err = c.ShouldBindJSON(&update)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	orderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = orderClient.UpdateCartItem(context.WithValue(context.Background(), "ginContext", c), &proto.CartItemRequest{
		UserId:  claims.UserID,
		GoodsId: cr.Id,
		Nums:    update.Num,
		Checked: update.Checked,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}
