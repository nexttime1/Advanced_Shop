package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"strconv"
)

func (gc *goodsController) GetBannerListView(c *gin.Context) error {
	ctx := c.Request.Context()
	list, err := gc.srv.Goods().BannerList(ctx, &empty.Empty{})
	if err != nil {
		return err

	}
	var response []good.BannerListResponse

	for _, model := range list.Data {
		response = append(response, good.BannerListResponse{
			Id:    model.Id,
			Index: model.Index,
			Image: model.Image,
			Url:   model.Url,
		})
	}

	common.OkWithList(c, response, list.Total)
	return nil
}

func (gc *goodsController) CreateBannerView(c *gin.Context) error {

	var cr good.BannerCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	ctx := c.Request.Context()
	bannerInfo, err := gc.srv.Goods().CreateBanner(ctx, &proto.BannerRequest{
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		return err

	}

	RMap := map[string]interface{}{
		"id": bannerInfo.Id,
	}
	common.OkWithData(c, RMap)
	return nil
}

func (gc *goodsController) DeleteBannerView(c *gin.Context) error {

	var cr good.BannerIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	_, err = gc.srv.Goods().DeleteBanner(ctx, &proto.BannerRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil
}

func (gc *goodsController) UpdateBannerView(c *gin.Context) error {

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	var cr good.BannerUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	_, err = gc.srv.Goods().UpdateBanner(ctx, &proto.BannerRequest{
		Id:    int32(id),
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "更新成功")
	return nil
}
