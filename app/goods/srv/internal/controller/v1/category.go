package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/gorm"
	"Advanced_Shop/pkg/log"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func convertCategoryDOToInfo(do *do.CategoryDO) *proto.CategoryInfoResponse {
	if do == nil {
		return nil
	}
	return &proto.CategoryInfoResponse{
		Id:               do.ID,
		Name:             do.Name,
		ParentCategoryID: do.ParentCategoryID,
		Level:            do.Level,
		IsTab:            do.IsTab,
	}
}

// GetAllCategorysList 获取所有分类列表
func (gs *goodsServer) GetAllCategorysList(ctx context.Context, empty *emptypb.Empty) (*proto.CategoryListResponse, error) {
	log.Info("GetAllCategory Call")
	// 1. 调用service层获取所有分类（ 默认按ID升序）
	categoryList, err := gs.srv.Category().ListAll(ctx, []string{"id asc"})
	if err != nil {
		log.Errorf("get all category list error: %v", err.Error())
		return nil, err
	}

	ret := proto.CategoryListResponse{
		JsonData: categoryList.JsonData,
		Total:    int32(categoryList.TotalCount),
	}

	return &ret, nil
}

// GetSubCategory 获取子分类列表
func (gs *goodsServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 1. 调用data层获取分类（已预加载三级）
	categoryDO, err := gs.srv.Category().Get(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("获取子分类失败，父分类ID: %d，错误: %v", request.Id, err)
		return nil, err
	}

	// 2. 构建返回结构体
	ret := &proto.SubCategoryListResponse{
		Info:  convertCategoryDOToInfo(categoryDO), // 根分类信息
		Total: int32(len(categoryDO.SubCategory)),  // 根分类的直接子分类
	}

	// 3. 转换直接子分类（二级），并自动携带三级分类（通过辅助函数递归填充）
	var subCategorys []*proto.CategoryInfoResponse
	for _, subDO := range categoryDO.SubCategory {
		subCategorys = append(subCategorys, convertCategoryDOToInfo(subDO))
	}
	ret.SubCategorys = subCategorys

	return ret, nil
}

// CreateCategory 创建分类
func (gs *goodsServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	// 1. Proto请求转换为DO
	categoryDO := &do.CategoryDO{
		Name:             request.Name,
		ParentCategoryID: request.ParentCategoryID,
		Level:            request.Level,
	}
	if request.IsTab != nil {
		categoryDO.IsTab = *request.IsTab
	}

	// 2. 调用service层创建方法
	err := gs.srv.Category().Create(ctx, categoryDO)
	if err != nil {
		log.Errorf("create category error: %v", err.Error())
		return nil, err
	}

	// 3. 转换为Proto响应
	return &proto.CategoryInfoResponse{
		Id:               categoryDO.ID,
		Name:             categoryDO.Name,
		ParentCategoryID: categoryDO.ParentCategoryID,
		Level:            int32(categoryDO.Level),
		IsTab:            categoryDO.IsTab,
	}, nil
}

// DeleteCategory 删除分类
func (gs *goodsServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	// 1. 调用service层删除方法
	err := gs.srv.Category().Delete(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("delete category error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 2. 返回空响应
	return &emptypb.Empty{}, nil
}

// UpdateCategory 更新分类
func (gs *goodsServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	// 1. Proto请求转换为DO
	categoryDO := &do.CategoryDO{
		Model:            gorm.Model{ID: request.Id},
		Name:             request.Name,
		ParentCategoryID: request.ParentCategoryID,
		Level:            request.Level,
	}
	if request.IsTab != nil {
		categoryDO.IsTab = *request.IsTab
	}

	// 2. 调用service层更新方法
	err := gs.srv.Category().Update(ctx, categoryDO)
	if err != nil {
		log.Errorf("update category error, id: %d, err: %v", request.Id, err.Error())
		return nil, err
	}

	// 3. 返回空响应
	return &emptypb.Empty{}, nil
}
