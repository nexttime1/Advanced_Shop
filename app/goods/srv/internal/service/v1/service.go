package v1

import (
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	v12 "Advanced_Shop/app/goods/srv/internal/data_search/v1"
)

type ServiceFactory interface {
	Goods() GoodsSrv
	Brands() BrandsSrv
	Category() CategorySrv
	CategoryBrands() CategoryBrandSrv
	Banner() BannerSrv
}

type serviceFactory struct {
	data       v1.DataFactory
	dataSearch v12.SearchFactory
}

func NewService(store v1.DataFactory, dataSearch v12.SearchFactory) ServiceFactory {
	return &serviceFactory{data: store, dataSearch: dataSearch}
}

var _ ServiceFactory = &serviceFactory{}

func (s *serviceFactory) Goods() GoodsSrv {
	return newGoods(s)
}

func (s *serviceFactory) Brands() BrandsSrv {
	return newBrand(s)
}
func (s *serviceFactory) Category() CategorySrv {
	return newCategory(s)
}

func (s *serviceFactory) Banner() BannerSrv {
	return newBanner(s)
}

func (s *serviceFactory) CategoryBrands() CategoryBrandSrv {
	return newCategoryBrand(s)
}
