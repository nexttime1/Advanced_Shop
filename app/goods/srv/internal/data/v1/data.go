package v1

import "gorm.io/gorm"

type DataFactory interface {
	Goods() GoodsStore
	Categorys() CategoryStore
	Brands() BrandsStore
	Banners() BannerStore
	CategoryBrands() GoodsCategoryBrandStore

	Begin() *gorm.DB
}
