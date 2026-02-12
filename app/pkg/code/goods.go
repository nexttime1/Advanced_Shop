package code

const (
	// ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 100501

	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound

	// ErrBrandNotFound - 500: Es unmarshal error.
	ErrEsUnmarshal

	// ErrBannerNotFound - 404: Banner not found.
	ErrBannerNotFound

	// ErrBrandNotFound - 404: Brand not found.
	ErrBrandNotFound

	// ErrCategoryBrandNotFound - 404: CategoryBrand not found.
	ErrCategoryBrandNotFound

	// ErrGoodsImageNotFound - 404: GoodImage not found.
	ErrGoodsImageNotFound
)
