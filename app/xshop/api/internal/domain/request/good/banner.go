package good

type BannerCreateRequest struct {
	Image string `form:"image" json:"image" binding:"url"`
	Index int32  `form:"index" json:"index" binding:"required"`
	Url   string `form:"url" json:"url" binding:"url"`
}

type BannerIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type BannerUpdateRequest struct {
	Image string `form:"image" json:"image"`
	Index int32  `form:"index" json:"index"`
	Url   string `form:"url" json:"url"`
}

type BannerListResponse struct {
	Id    int32  `json:"id,omitempty"`    // 主键ID
	Index int32  `json:"index,omitempty"` // 排序优先级(值越小越靠前)
	Image string `json:"image,omitempty"` // 图片链接地址
	Url   string `json:"url,omitempty"`   // 点击跳转链接
}
