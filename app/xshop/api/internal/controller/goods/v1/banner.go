package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"Advanced_Shop/pkg/common/core"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"strconv"
)

func (gc *goodsController) GetBannerListView(c *gin.Context) {

	list, err := gc.srv.Goods().BannerList(c, &empty.Empty{})
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
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

	core.OkWithList(c, response, list.Total)

}

func (gc *goodsController) CreateBannerView(c *gin.Context) {

	var cr good.BannerCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	bannerInfo, err := gc.srv.Goods().CreateBanner(c, &proto.BannerRequest{
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	RMap := map[string]interface{}{
		"id": bannerInfo.Id,
	}
	core.OkWithData(c, RMap)
}

func (gc *goodsController) DeleteBannerView(c *gin.Context) {

	var cr good.BannerIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().DeleteBanner(c, &proto.BannerRequest{
		Id: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}

func (gc *goodsController) UpdateBannerView(c *gin.Context) {

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	var cr good.BannerUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().UpdateBanner(c, &proto.BannerRequest{
		Id:    int32(id),
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")

}
