package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"Advanced_Shop/app/xshop/api/internal/service"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/log"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"strconv"
)

type goodsController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewGoodsController(srv service.ServiceFactory, trans ut.Translator) *goodsController {
	return &goodsController{
		srv:   srv,
		trans: trans,
	}
}

func (gc *goodsController) GetGoodListView(c *gin.Context) {
	log.Info("goods list function called ...")

	var cr good.GoodListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	list, err := gc.srv.Goods().List(c, &proto.GoodsFilterRequest{
		PriceMin:      cr.PriceMin,
		PriceMax:      cr.PriceMax,
		IsHot:         cr.IsHot,
		IsNew:         cr.IsNew,
		TopCategoryID: cr.TopCategoryID,
		Pages:         cr.Page,
		PagePerNums:   cr.Limit,
		KeyWords:      cr.Key,
		BrandID:       cr.BrandID,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	var response []good.GoodsInfoResponse
	for _, model := range list.Data {
		info := good.GoodsInfoResponse{
			ID:              model.Id,
			CategoryID:      model.CategoryId,
			Name:            model.Name,
			GoodsSn:         model.GoodsSn,
			ClickNum:        model.ClickNum,
			SoldNum:         model.SoldNum,
			FavNum:          model.FavNum,
			Stocks:          model.Stocks,
			MarketPrice:     model.MarketPrice,
			ShopPrice:       model.ShopPrice,
			GoodsBrief:      model.GoodsBrief,
			GoodsDesc:       model.GoodsDesc,
			ShipFree:        model.ShipFree,
			Images:          model.Images,
			DescImages:      model.DescImages,
			GoodsFrontImage: model.GoodsFrontImage,
			IsNew:           model.IsNew,
			IsHot:           model.IsHot,
			OnSale:          model.OnSale,
			AddTime:         model.AddTime,
			Category: good.CategoryBriefInfoResponse{
				ID:   model.Category.Id,
				Name: model.Category.Name,
			},
			Brand: good.BrandInfoResponse{
				ID:   model.Brand.Id,
				Name: model.Brand.Name,
				Logo: model.Brand.Logo,
			},
		}
		response = append(response, info)
	}
	core.OkWithList(c, response, list.Total)

}

func (gc *goodsController) CreateGoodView(c *gin.Context) {

	var cr good.GoodCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	_, err = gc.srv.Goods().CreateGoods(c, &proto.CreateGoodsInfo{
		Name:            cr.Name,
		GoodsSn:         cr.GoodsSn,
		Stocks:          cr.Stocks,
		MarketPrice:     cr.MarketPrice,
		ShopPrice:       cr.ShopPrice,
		GoodsBrief:      cr.GoodsBrief,
		ShipFree:        cr.ShipFree,
		Images:          cr.Images,
		DescImages:      cr.DescImages,
		GoodsFrontImage: cr.FrontImage,
		CategoryId:      cr.CategoryId,
		BrandId:         cr.Brand,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "创建成功")
}

func (gc *goodsController) GoodDetailView(c *gin.Context) {

	var cr good.GoodDetailRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	goodInfo, err := gc.srv.Goods().GetGoodsDetail(c, &proto.GoodInfoRequest{
		Id: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	// TODO
	//stockClient, clientConn, err := connect.StockConnectService(c)
	//if err != nil {
	//	return
	//}
	//defer clientConn.Close()
	//detail, err := stockClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
	//	GoodsId: goodInfo.Id,
	//})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	response := good.GoodsInfoResponse{
		ID:              goodInfo.Id,
		CategoryID:      goodInfo.CategoryId,
		Name:            goodInfo.Name,
		GoodsSn:         goodInfo.GoodsSn,
		ClickNum:        goodInfo.ClickNum,
		SoldNum:         goodInfo.SoldNum,
		FavNum:          goodInfo.FavNum,
		Stocks:          detail.Num,
		MarketPrice:     goodInfo.MarketPrice,
		ShopPrice:       goodInfo.ShopPrice,
		GoodsBrief:      goodInfo.GoodsBrief,
		GoodsDesc:       goodInfo.GoodsDesc,
		ShipFree:        goodInfo.ShipFree,
		Images:          goodInfo.Images,
		DescImages:      goodInfo.DescImages,
		GoodsFrontImage: goodInfo.GoodsFrontImage,
		IsNew:           goodInfo.IsNew,
		IsHot:           goodInfo.IsHot,
		OnSale:          goodInfo.OnSale,
		AddTime:         goodInfo.AddTime,
		Category: good.CategoryBriefInfoResponse{
			ID:   goodInfo.Category.Id,
			Name: goodInfo.Category.Name,
		},
		Brand: good.BrandInfoResponse{
			ID:   goodInfo.Brand.Id,
			Name: goodInfo.Brand.Name,
			Logo: goodInfo.Brand.Logo,
		},
	}

	core.OkWithData(c, response)

}

func (gc *goodsController) GoodUpdateView(c *gin.Context) {

	var cr good.GoodUpdateRequest
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	_, err = gc.srv.Goods().UpdateGoods(c, &proto.CreateGoodsInfo{
		Id:              int32(id),
		Name:            cr.Name,
		GoodsSn:         cr.GoodsSn,
		Stocks:          cr.Stocks,
		MarketPrice:     cr.MarketPrice,
		ShopPrice:       cr.ShopPrice,
		GoodsBrief:      cr.GoodsBrief,
		ShipFree:        cr.ShipFree,
		Images:          cr.Images,
		DescImages:      cr.DescImages,
		GoodsFrontImage: cr.FrontImage,
		CategoryId:      cr.CategoryId,
		BrandId:         cr.Brand,
	})

	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")

}

func (gc *goodsController) GoodPatchUpdateView(c *gin.Context) {

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	var cr good.GoodPatchUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().UpdateGoods(c, &proto.CreateGoodsInfo{
		Id:         int32(id),
		IsNew:      cr.IsNew,
		IsHot:      cr.IsHot,
		OnSale:     cr.OnSale,
		CategoryId: 0,
		BrandId:    0,
	})

	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")
}

func (gc *goodsController) GoodDeleteView(c *gin.Context) {

	var cr good.GoodDeleteRequest

	err := c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().DeleteGoods(c, &proto.DeleteGoodsInfo{
		Id: cr.Id,
	})

	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}
