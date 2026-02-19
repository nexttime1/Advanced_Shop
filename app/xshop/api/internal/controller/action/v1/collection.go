package v1

import (
	p1 "Advanced_Shop/api/action/v1"
	p2 "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/action"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ac *actionController) CollectionListView(c *gin.Context) {
	log.Info("collection list function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	ctx := c.Request.Context()
	list, err := ac.srv.Collection().GetFavList(ctx, &p1.UserFavRequest{
		UserId: userID,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	if list.Total == 0 {
		core.OkWithList(c, list.Data, 0)
		return
	}
	var idList []int32
	for _, model := range list.Data {
		idList = append(idList, model.GoodsId)
	}

	goodsInfo, err := ac.srv.Goods().BatchGetGoods(ctx, &p2.BatchGoodsIdInfo{
		Id: idList,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
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
	core.OkWithList(c, response, goodsInfo.Total)

}

func (ac *actionController) CollectionAddView(c *gin.Context) {
	log.Info("collection add function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr action.CollectionAddRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}
	ctx := c.Request.Context()
	_, err = ac.srv.Collection().AddUserFav(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	core.OkWithMessage(c, "收藏成功")
}

func (ac *actionController) CollectionDeleteView(c *gin.Context) {
	log.Info("collection delete function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	var cr action.CollectionIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}

	ctx := c.Request.Context()
	_, err = ac.srv.Collection().DeleteUserFav(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}

func (ac *actionController) CollectionDetailView(c *gin.Context) {
	log.Info("collection delete function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	var cr action.CollectionIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}

	ctx := c.Request.Context()
	_, err = ac.srv.Collection().GetUserFavDetail(ctx, &p1.UserFavRequest{
		UserId:  userID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "存在")

}
