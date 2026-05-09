package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (gc *goodsController) BrandListView(c *gin.Context) error {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	list, err := gc.srv.Goods().BrandList(c, &proto.BrandFilterRequest{
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})
	if err != nil {
		return err
	}

	common.OkWithList(c, list.Data, list.Total)
	return nil
}

func (gc *goodsController) CreateBrandView(c *gin.Context) error {
	var cr good.BrandCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	brandInfo, err := gc.srv.Goods().CreateBrand(ctx, &proto.BrandRequest{
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		return err
	}
	RMap := map[string]interface{}{
		"id": brandInfo.Id,
	}
	common.OkWithData(c, RMap)
	return nil
}

func (gc *goodsController) UpdateBrandView(c *gin.Context) error {
	var cr good.BrandUpdateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	_, err = gc.srv.Goods().UpdateBrand(ctx, &proto.BrandRequest{
		Id:   int32(id),
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		return err
	}

	common.OkWithMessage(c, "修改成功")
	return nil
}

func (gc *goodsController) DeleteBrandView(c *gin.Context) error {
	var cr good.BrandIdRequest

	if err := c.ShouldBindUri(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	ctx := c.Request.Context()
	_, err := gc.srv.Goods().DeleteBrand(ctx, &proto.BrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil
}

//第三张表

func (gc *goodsController) CategoryBrandListView(c *gin.Context) error {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	list, err := gc.srv.Goods().CategoryBrandList(ctx, &proto.CategoryBrandFilterRequest{
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})
	if err != nil {
		return err
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

	common.OkWithList(c, response, list.Total)
	return nil
}

func (gc *goodsController) CategoryAllBrandView(c *gin.Context) error {
	var cr good.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	list, err := gc.srv.Goods().GetCategoryBrandList(ctx, &proto.CategoryInfoRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithList(c, list.Data, list.Total)
	return nil
}

func (gc *goodsController) CreateCategoryBrandView(c *gin.Context) error {
	var cr good.CreateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	Info, err := gc.srv.Goods().CreateCategoryBrand(ctx, &proto.CategoryBrandRequest{
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		return err
	}
	RMap := map[string]interface{}{
		"id": Info.Id,
	}
	common.OkWithData(c, RMap)
	return nil
}

func (gc *goodsController) DeleteCategoryBrandView(c *gin.Context) error {
	var cr good.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	_, err := gc.srv.Goods().DeleteCategoryBrand(ctx, &proto.CategoryBrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil

}

func (gc *goodsController) UpdateCategoryBrandView(c *gin.Context) error {
	var cr good.UpdateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	_, err = gc.srv.Goods().UpdateCategoryBrand(ctx, &proto.CategoryBrandRequest{
		Id:         int32(id),
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "更新成功")
	return nil

}
