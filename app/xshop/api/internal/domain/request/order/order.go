package order

type OrderCreateRequest struct {
	Post    string `json:"post" binding:"required"`
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required,mobile"`
}
type OrderCreateResponse struct {
	OrderSn   string `json:"order_sn"`
	AlipayUrl string `json:"alipay_url"`
}

type OrderIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type OrderUpdateRequest struct {
	Image string `form:"image" json:"image"`
	Index int32  `form:"index" json:"index"`
	Url   string `form:"url" json:"url"`
}

type OrderListResponse struct {
	Id      int32   `json:"id"`
	UserId  int32   `json:"user_id"`
	OrderSn string  `json:"order_sn"`
	PayType string  `json:"pay_type"`
	Status  string  `json:"status"`
	Post    string  `json:"post"`
	Total   float32 `json:"total"`
	Address string  `json:"address"`
	Name    string  `json:"name"`
	Mobile  string  `json:"mobile"`
}

type OrderDetailResponse struct {
	Id        int32      `json:"id"`
	UserId    int32      `json:"user_id"`
	OrderSn   string     `json:"order_sn"`
	PayType   string     `json:"pay_type"`
	Status    string     `json:"status"`
	Post      string     `json:"post"`
	Total     float32    `json:"total"`
	Address   string     `json:"address"`
	Name      string     `json:"name"`
	Mobile    string     `json:"mobile"`
	GoodInfo  []GoodInfo `json:"goods"`
	AlipayUrl string     `json:"alipay_url"`
}

type GoodInfo struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`

	Image string  `json:"image"`
	Price float32 `json:"price"`
	Nums  int32   `json:"nums"`
}
