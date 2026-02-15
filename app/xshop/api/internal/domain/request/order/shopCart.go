package order

type CartListResponse struct {
	Id          int32    `json:"id"`
	GoodID      int32    `json:"good_id"`
	Name        string   `form:"name" json:"name"`
	GoodsSn     string   `form:"goods_sn" json:"goods_sn"`
	Stocks      int32    `form:"stocks" json:"stocks"` // 有值时校验 ≥1
	CategoryId  int32    `form:"category"`
	MarketPrice float32  `form:"market_price"`
	GoodPrice   float32  `form:"good_price" json:"good_price"`
	GoodsBrief  string   `form:"goods_brief" json:"goods_brief"`
	Images      []string `form:"images" json:"images"`
	DescImages  []string `form:"desc_images" json:"desc_images"`
	ShipFree    *bool    `form:"ship_free" json:"ship_free"`
	FrontImage  string   `form:"front_image" json:"front_image"`
	Brand       int32    `form:"brand" json:"brand"`
	Chacked     *bool    `form:"chacked" json:"chacked"`
}

type CartAddRequest struct {
	GoodID int32 `json:"good_id" binding:"required"`
	Num    int32 `json:"num" binding:"required" min:"1"`
}

type CartAddResponse struct {
	Id int32 `json:"id"`
}

type CartIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type CartUpdateRequest struct {
	Num     int32 `json:"num" binding:"required" min:"1"`
	Checked *bool `json:"checked"`
}
