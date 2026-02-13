package request

import "Advanced_Shop/app/pkg/common"

type GoodsFilter struct {
	PriceMin    int32  `form:"pmin"`
	PriceMax    int32  `form:"pmax"`
	IsHot       bool   `form:"ih"`
	IsNew       bool   `form:"in"`
	IsTab       bool   `form:"it"`
	TopCategory int32  `form:"c"`
	Pages       int32  `form:"p"`
	PagePerNums int32  `form:"pnum"`
	KeyWords    string `form:"q"`
	Brand       int32  `form:"b"`
}

type GoodListRequest struct {
	common.PageInfo
	IsHot         bool  `form:"is_hot"`
	IsNew         bool  `form:"is_new"`
	PriceMax      int32 `form:"price_max"`
	PriceMin      int32 `form:"price_min"`
	BrandID       int32 `form:"brand_id"`
	TopCategoryID int32 `form:"top_category_id"`
}

type GoodCreateRequest struct {
	Name        string   `form:"name" json:"name" binding:"required,min=2,max=100"`
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"required,min=2,lt=20"`
	Stocks      int32    `form:"stocks" json:"stocks" binding:"required,min=1"`
	CategoryId  int32    `form:"category" json:"category" binding:"required"`
	MarketPrice float32  `form:"market_price" json:"market_price" binding:"required"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"required,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"required,min=3"`
	Images      []string `form:"images" json:"images" binding:"required,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images" binding:"required"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free" binding:"required"`
	FrontImage  string   `form:"front_image" json:"front_image" binding:"required,url"`
	Brand       int32    `form:"brand" json:"brand" binding:"required"`
}

type GoodsInfoResponse struct {
	ID              int32                     `json:"id"`                  // 商品ID
	CategoryID      int32                     `json:"category_id"`         // 分类ID
	Name            string                    `json:"name"`                // 商品名称
	GoodsSn         string                    `json:"goods_sn"`            // 商品编号
	ClickNum        int32                     `json:"click_num"`           // 点击数
	SoldNum         int32                     `json:"sold_num"`            // 销量
	FavNum          int32                     `json:"fav_num"`             // 收藏数
	Stocks          int32                     `json:"stocks"`              // 库存
	MarketPrice     float32                   `json:"market_price"`        // 市场价
	ShopPrice       float32                   `json:"shop_price"`          // 店铺价
	GoodsBrief      string                    `json:"goods_brief"`         // 商品简介
	GoodsDesc       string                    `json:"goods_desc"`          // 商品详情
	ShipFree        *bool                     `json:"ship_free,omitempty"` // 是否包邮（optional，指针表示可选）
	Images          []string                  `json:"images"`              // 商品图片（repeated）
	DescImages      []string                  `json:"desc_images"`         // 详情图片（repeated）
	GoodsFrontImage string                    `json:"goods_front_image"`   // 商品封面图
	IsNew           *bool                     `json:"is_new,omitempty"`    // 是否新品（optional）
	IsHot           *bool                     `json:"is_hot,omitempty"`    // 是否热门（optional）
	OnSale          *bool                     `json:"on_sale,omitempty"`   // 是否上架（optional）
	AddTime         int64                     `json:"add_time"`            // 添加时间
	Category        CategoryBriefInfoResponse `json:"category"`            // 分类信息
	Brand           BrandInfoResponse         `json:"brand"`               // 品牌信息
}

// CategoryBriefInfoResponse 对应 Protobuf 的 CategoryBriefInfoResponse 消息
type CategoryBriefInfoResponse struct {
	ID   int32  `json:"id"`   // 分类ID
	Name string `json:"name"` // 分类名称
}

// BrandInfoResponse 对应 Protobuf 的 BrandInfoResponse 消息
type BrandInfoResponse struct {
	ID   int32  `json:"id"`   // 品牌ID
	Name string `json:"name"` // 品牌名称
	Logo string `json:"logo"` // 品牌logo
}

type GoodDetailRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type GoodUpdateRequest struct {
	Name        string   `form:"name" json:"name" binding:"omitempty,min=2,max=100"`       // 有值时校验长度 2-100，
	GoodsSn     string   `form:"goods_sn" json:"goods_sn" binding:"omitempty,min=2,lt=20"` // 有值时校验长度 2-19，
	Stocks      int32    `form:"stocks" json:"stocks" binding:"omitempty,min=1"`           // 有值时校验 ≥1
	CategoryId  int32    `form:"category" json:"category"`
	MarketPrice float32  `form:"market_price" json:"market_price"`
	ShopPrice   float32  `form:"shop_price" json:"shop_price" binding:"omitempty,min=0"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief" binding:"omitempty,min=3"`
	Images      []string `form:"images" json:"images" binding:"omitempty,min=1"`
	DescImages  []string `form:"desc_images" json:"desc_images"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free"`
	FrontImage  string   `form:"front_image" json:"front_image"`
	Brand       int32    `form:"brand" json:"brand"`
}

type GoodPatchUpdateRequest struct {
	IsNew  *bool `form:"is_new" json:"is_new"`
	IsHot  *bool `form:"is_hot" json:"is_hot"`
	OnSale *bool `form:"on_sale" json:"on_sale"`
}

type GoodDeleteRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}
