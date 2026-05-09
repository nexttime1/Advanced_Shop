package goods

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/good"
	"Advanced_Shop/pkg/log"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"strconv"
)

func (gc *goodsController) GetAllCategoryView(c *gin.Context) error {
	log.Info("GetAllCategory Call")
	ctx := c.Request.Context()
	list, err := gc.srv.Goods().GetAllCategorysList(ctx, &empty.Empty{})
	if err != nil {
		return err

	}
	var response []interface{}

	err = json.Unmarshal([]byte(list.JsonData), &response)
	if err != nil {
		return err
	}

	common.OkWithList(c, response, list.Total)
	return nil

}

func (gc *goodsController) GetSubCategoryView(c *gin.Context) error {
	var cr good.CategoryIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	ctx := c.Request.Context()
	// 调用gRPC服务获取数据（适配新的proto结构）
	categoryInfo, err := gc.srv.Goods().GetSubCategory(ctx, &proto.CategoryListRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}

	// 转换proto结构为Web层响应结构
	response := ProtoToWebSubCategory(categoryInfo)
	common.OkWithData(c, response)
	return nil
}

// 辅助函数：转换CategoryInfoResponse → CategoryInfoResponse
func protoToWebCategoryInfo(protoInfo *proto.CategoryInfoResponse) *good.CategoryInfoResponse {
	if protoInfo == nil {
		return nil
	}
	webInfo := &good.CategoryInfoResponse{
		Id:               protoInfo.Id,
		Name:             protoInfo.Name,
		ParentCategoryID: protoInfo.ParentCategoryID,
		Level:            protoInfo.Level,
		IsTab:            protoInfo.IsTab,
	}

	// 递归处理子分类（三级/四级）
	if protoInfo.SubCategorys != nil && len(protoInfo.SubCategorys) > 0 {
		webSubs := make([]*good.CategoryInfoResponse, 0, len(protoInfo.SubCategorys))
		for _, protoSub := range protoInfo.SubCategorys {
			webSubs = append(webSubs, protoToWebCategoryInfo(protoSub))
		}
		webInfo.SubCategories = webSubs
	}
	return webInfo
}

// ProtoToWebSubCategory 转换SubCategoryListResponse（proto）→ SubCategoryResponse（web）
func ProtoToWebSubCategory(protoCategory *proto.SubCategoryListResponse) *good.SubCategoryResponse {
	if protoCategory == nil {
		return nil
	}

	// 转换根分类信息
	webInfo := protoToWebCategoryInfo(protoCategory.Info)

	// 转换直接子分类列表
	var webSubs []*good.CategoryInfoResponse
	if protoCategory.SubCategorys != nil && len(protoCategory.SubCategorys) > 0 {
		webSubs = make([]*good.CategoryInfoResponse, 0, len(protoCategory.SubCategorys))
		for _, protoSub := range protoCategory.SubCategorys {
			webSubs = append(webSubs, protoToWebCategoryInfo(protoSub))
		}
	}

	// 3. 构建Web层最终响应
	webCategory := &good.SubCategoryResponse{
		Total:         protoCategory.Total, // 直接子分类数量
		Info:          webInfo,             // 根分类信息（含嵌套子分类）
		SubCategories: webSubs,             // 直接子分类列表（含三级）
	}
	return webCategory
}

func (gc *goodsController) CreateCategoryView(c *gin.Context) error {

	var cr good.CategoryCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	ctx := c.Request.Context()
	category, err := gc.srv.Goods().CreateCategory(ctx, &proto.CategoryInfoRequest{
		Name:             cr.Name,
		ParentCategoryID: cr.ParentCategory,
		Level:            cr.Level,
		IsTab:            cr.IsTab,
	})
	if err != nil {
		return err
	}
	RMap := map[string]interface{}{
		"id": category.Id,
	}
	common.OkWithData(c, RMap)
	return nil

}

func (gc *goodsController) UpdateCategoryView(c *gin.Context) error {

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}

	var cr good.UpdateCategoryRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)
	}
	ctx := c.Request.Context()
	_, err = gc.srv.Goods().UpdateCategory(ctx, &proto.CategoryInfoRequest{
		Id:    int32(id),
		Name:  cr.Name,
		IsTab: cr.IsTab,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "更新成功")
	return nil

}

func (gc *goodsController) DeleteCategoryView(c *gin.Context) error {

	var cr good.CategoryIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, gc.trans)

	}
	ctx := c.Request.Context()
	_, err = gc.srv.Goods().DeleteCategory(ctx, &proto.DeleteCategoryRequest{
		Id: cr.Id,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil
}
