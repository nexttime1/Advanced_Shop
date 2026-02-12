package service

type BannerUpdateServiceMap struct {
	Image string `structs:"image"`
	Url   string `structs:"url"`
	Index int32  `structs:"index"`
}

type BrandUpdateServiceMap struct {
	Name string `structs:"name"`
	Logo string `structs:"logo"`
}

type CategoryBrandUpdateServiceMap struct {
	CategoryID int32 `structs:"category_id"` // 字段名需和gorm标签/数据库列名匹配
	BrandsID   int32 `structs:"brands_id"`
}

type GoodUpdateServiceMap struct {
	Name        string  `structs:"name"`
	GoodsSn     string  `structs:"goods_sn"`
	Stocks      int32   `structs:"stocks"`
	MarketPrice float32 `structs:"market_price"`
	ShopPrice   float32 `structs:"shop_price"`
	GoodsBrief  string  `structs:"goods_brief"`
	ShipFree    *bool   `structs:"ship_free"`
	IsNew       *bool   `structs:"is_new"`
	IsHot       *bool   `structs:"is_hot"`
	OnSale      *bool   `structs:"on_sale"`
	CategoryId  int32   `structs:"category_id"`
	Brand       int32   `structs:"brands_id"`
}
