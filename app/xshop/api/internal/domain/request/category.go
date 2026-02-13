package request

type CategoryIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type CategoryCreateRequest struct {
	Name           string `form:"name" json:"name" binding:"required,min=3,max=20"`
	ParentCategory int32  `form:"parent" json:"parent"`
	Level          int32  `form:"level" json:"level" binding:"required,oneof=1 2 3"`
	IsTab          *bool  `form:"is_tab" json:"is_tab" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name  string `form:"name" json:"name" binding:"required,min=3,max=20"`
	IsTab *bool  `form:"is_tab" json:"is_tab"`
}

type SubCategoryResponse struct {
	Total         int32                   `json:"total"`          // 直接子分类数量
	Info          *CategoryInfoResponse   `json:"info"`           // 当前查询的根分类信息
	SubCategories []*CategoryInfoResponse `json:"sub_categories"` // 根分类的直接子分类列表（二级）
}

// CategoryInfoResponse Web层分类基础信息（支持嵌套子分类）
type CategoryInfoResponse struct {
	Id               int32                   `json:"id"`                 // 分类ID
	Name             string                  `json:"name"`               // 分类名称
	ParentCategoryID int32                   `json:"parent_category_id"` // 父分类ID（对应proto的parentCategoryID）
	Level            int32                   `json:"level"`              // 分类层级
	IsTab            bool                    `json:"is_tab"`             // 是否为Tab
	SubCategories    []*CategoryInfoResponse `json:"sub_categories"`     // 子分类列表（三级/四级）
}
