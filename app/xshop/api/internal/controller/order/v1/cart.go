package v1

import (
	proto "Advanced_Shop/api/order/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/order"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (oc orderController) CartListView(c *gin.Context) error {
	log.Info("Cart list function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	ctx := c.Request.Context()
	response, total, err := oc.srv.Order().CartItemList(ctx, &proto.UserInfo{
		Id: userID,
	})

	if err != nil {
		return err
	}

	common.OkWithList(c, response, total)
	return nil

}

func (oc orderController) DeleteCartItemView(c *gin.Context) error {
	log.Info("DeleteCartItemView function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr order.CartIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, nil)

	}
	ctx := c.Request.Context()
	_, err = oc.srv.Order().DeleteCartItem(ctx, &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil
}

func (oc orderController) AddItemView(c *gin.Context) error {
	log.Info("AddItemView function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr order.CartAddRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, nil)

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
	common.OkWithData(c, response.Id)
	return nil
}

func (oc orderController) UpdatePatchView(c *gin.Context) error {
	log.Info("UpdatePatchView function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr order.CartIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, nil)

	}
	var update order.CartUpdateRequest
	err = c.ShouldBindJSON(&update)
	if err != nil {
		return gin2.HandleValidatorError(c, err, nil)

	}
	ctx := c.Request.Context()
	_, err = oc.srv.Order().UpdateCartItem(ctx, &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: cr.Id,
		Nums:    update.Num,
		Checked: update.Checked,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "更新成功")
	return nil

}
