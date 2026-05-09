package v1

import (
	p1 "Advanced_Shop/api/action/v1"
	p2 "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/action"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ac *actionController) CollectionListView(c *gin.Context) error {
	log.Info("collection list function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	ctx := c.Request.Context()
	list, err := ac.srv.Collection().GetFavList(ctx, &p1.UserFavRequest{
		UserId: userID,
	})
	if err != nil {
		return err
	}
	if list.Total == 0 {
		common.OkWithList(c, list.Data, 0)
		return nil
	}
	var idList []int32
	for _, model := range list.Data {
		idList = append(idList, model.GoodsId)
	}

	goodsInfo, err := ac.srv.Goods().BatchGetGoods(ctx, &p2.BatchGoodsIdInfo{
		Id: idList,
	})
	if err != nil {
		return err
	}
	var response []action.CollectionListResponse

	for _, Info := range list.Data {
		for _, goodModel := range goodsInfo.Data {
			if Info.GoodsId == goodModel.Id {
				response = append(response, action.CollectionListResponse{
					GoodId:    goodModel.Id,
					Name:      goodModel.Name,
					ShopPrice: goodModel.ShopPrice,
				})
				break
			}
		}
	}
	common.OkWithList(c, response, goodsInfo.Total)
	return nil
}

func (ac *actionController) CollectionAddView(c *gin.Context) error {
	log.Info("collection add function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr action.CollectionAddRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

	}
	ctx := c.Request.Context()
	_, err = ac.srv.Collection().AddUserFav(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		return err
	}

	common.OkWithMessage(c, "收藏成功")
	return nil
}

func (ac *actionController) CollectionDeleteView(c *gin.Context) error {
	log.Info("collection delete function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	var cr action.CollectionIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

	}

	ctx := c.Request.Context()
	_, err = ac.srv.Collection().DeleteUserFav(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")

	return nil
}

func (ac *actionController) CollectionDetailView(c *gin.Context) error {
	log.Info("collection delete function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	var cr action.CollectionIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

	}

	ctx := c.Request.Context()
	_, err = ac.srv.Collection().GetUserFavDetail(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "存在")
	return nil
}
