//go:generate codegen -type=int

package code

// Goods: goods service errors.
// Code must start with 1005xx.
const (
	// ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 100501

	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound

	// ErrEsUnmarshal - 500: Elasticsearch unmarshal error.
	ErrEsUnmarshal

	// ErrBannerNotFound - 404: Banner not found.
	ErrBannerNotFound

	// ErrBrandNotFound - 404: Brand not found.
	ErrBrandNotFound

	// ErrCategoryBrandNotFound - 404: CategoryBrand not found.
	ErrCategoryBrandNotFound

	// ErrGoodsImageNotFound - 404: GoodsImage not found.
	ErrGoodsImageNotFound

	// ErrJsonUnmarshal - 500: JSON unmarshal error.
	ErrJsonUnmarshal
)
