package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"Advanced_Shop/pkg/common/core"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (gc *goodsController) BrandListView(c *gin.Context) {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	list, err := gc.srv.Goods().BrandList(c, &proto.BrandFilterRequest{
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	core.OkWithList(c, list.Data, list.Total)

}

func (gc *goodsController) CreateBrandView(c *gin.Context) {
	var cr good.BrandCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	brandInfo, err := gc.srv.Goods().CreateBrand(c, &proto.BrandRequest{
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	RMap := map[string]interface{}{
		"id": brandInfo.Id,
	}
	core.OkWithData(c, RMap)

}

func (gc *goodsController) UpdateBrandView(c *gin.Context) {
	var cr good.BrandUpdateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().UpdateBrand(c, &proto.BrandRequest{
		Id:   int32(id),
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}

	core.OkWithMessage(c, "修改成功")

}

func (gc *goodsController) DeleteBrandView(c *gin.Context) {
	var cr good.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	_, err := gc.srv.Goods().DeleteBrand(c, &proto.BrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}

//第三张表

func (gc *goodsController) CategoryBrandListView(c *gin.Context) {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	list, err := gc.srv.Goods().CategoryBrandList(c, &proto.CategoryBrandFilterRequest{
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	var response []good.BrandCategoryItem
	for _, model := range list.Data {
		response = append(response, good.BrandCategoryItem{
			Brand: good.Brand{
				Id:   model.Brand.Id,
				Name: model.Brand.Name,
				Logo: model.Brand.Logo,
			},

			Category: good.Category{
				Id:               model.Category.Id,
				Name:             model.Category.Name,
				ParentCategoryID: model.Category.ParentCategoryID,
				Level:            model.Category.Level,
				IsTab:            model.Category.IsTab,
			},
		})
	}

	core.OkWithList(c, response, list.Total)

}

func (gc *goodsController) CategoryAllBrandView(c *gin.Context) {
	var cr good.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	list, err := gc.srv.Goods().GetCategoryBrandList(c, &proto.CategoryInfoRequest{
		Id: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithList(c, list.Data, list.Total)

}

func (gc *goodsController) CreateCategoryBrandView(c *gin.Context) {
	var cr good.CreateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	Info, err := gc.srv.Goods().CreateCategoryBrand(c, &proto.CategoryBrandRequest{
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	RMap := map[string]interface{}{
		"id": Info.Id,
	}
	core.OkWithData(c, RMap)
}

func (gc *goodsController) DeleteCategoryBrandView(c *gin.Context) {
	var cr good.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err := gc.srv.Goods().DeleteCategoryBrand(c, &proto.CategoryBrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}

func (gc *goodsController) UpdateCategoryBrandView(c *gin.Context) {
	var cr good.UpdateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		gin2.HandleValidatorError(c, err, gc.trans)
		return
	}

	_, err = gc.srv.Goods().UpdateCategoryBrand(c, &proto.CategoryBrandRequest{
		Id:         int32(id),
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")

}
