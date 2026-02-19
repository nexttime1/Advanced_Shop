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
	res, err := oc.srv.Order().CreateCartItem(ctx, &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: cr.GoodID,
		Nums:    cr.Num,
	})
	response := order.CartAddResponse{
		Id: res.Id,
	}
	core.OkWithData(c, response.Id)

}

func (oc orderController) UpdatePatchView(c *gin.Context) {
	log.Info("UpdatePatchView function called ...")
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
	var update order.CartUpdateRequest
	err = c.ShouldBindJSON(&update)
	if err != nil {
		gin2.HandleValidatorError(c, err, nil)
		return
	}
	ctx := c.Request.Context()
	_, err = oc.srv.Order().UpdateCartItem(ctx, &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: cr.Id,
		Nums:    update.Num,
		Checked: update.Checked,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")

}
