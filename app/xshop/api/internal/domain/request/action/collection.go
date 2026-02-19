package action

type CollectionRequest struct {
	GoodID int32 `json:"good_id" binding:"required"`
}

type CollectionListResponse struct {
	GoodId    int32   `json:"good_id"`
	Name      string  `json:"name"`
	ShopPrice float32 `json:"shop_price"`
}

type CollectionAddRequest struct {
	GoodId int32 `json:"good_id" binding:"required,min=1"`
}

type CollectionIdRequest struct {
	GoodId int32 `uri:"good_id" binding:"required,min=1"`
}
