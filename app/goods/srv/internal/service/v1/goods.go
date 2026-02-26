package v1

import (
	proto "Advanced_Shop/api/goods/v1"
	v1 "Advanced_Shop/app/goods/srv/internal/data/v1"
	v12 "Advanced_Shop/app/goods/srv/internal/data_search/v1"
	"Advanced_Shop/app/goods/srv/internal/domain/do"
	"Advanced_Shop/app/goods/srv/internal/domain/dto"
	"context"
	"sync"

	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/log"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodsSrv interface {
	// List 商品列表
	List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error)

	// Get 商品详情
	Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error)

	// Create 创建商品
	Create(ctx context.Context, goods *v1.GoodsInfo) error

	// Update 更新商品
	Update(ctx context.Context, goods *v1.GoodsInfo) error

	// CreateInTxn 事务
	CreateInTxn(ctx context.Context, goods *v1.GoodsInfo) error
	// UpdateInTxn 事务
	UpdateInTxn(ctx context.Context, goods *v1.GoodsInfo) error

	// Delete 删除商品
	Delete(ctx context.Context, ID uint64) error

	// BatchGet 批量查询商品
	BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error)
}

type goodsService struct {
	//工厂
	data v1.DataFactory // 数据层（MySQL）

	searchData v12.SearchFactory // 搜索层（ES）

}

func newGoods(srv *serviceFactory) GoodsSrv {
	return &goodsService{
		data:       srv.data,
		searchData: srv.dataSearch,
	}
}

// 遍历树结构
func retrieveIDs(category *do.CategoryDO) []uint64 {
	ids := []uint64{}
	if category == nil || category.ID == 0 {
		return ids
	}
	ids = append(ids, uint64(category.ID))
	for _, child := range category.SubCategory {
		subids := retrieveIDs(child)
		ids = append(ids, subids...)
	}
	return ids
}

func (gs *goodsService) List(ctx context.Context, opts metav1.ListMeta, req *proto.GoodsFilterRequest, orderby []string) (*dto.GoodsDTOList, error) {
	searchReq := v12.GoodsFilterRequest{
		GoodsFilterRequest: req,
	}
	if req.TopCategoryID > 0 {
		category, err := gs.data.NewMysql().Categorys().Get(ctx, uint64(req.TopCategoryID))
		if err != nil {
			log.Errorf("categorydata.NewMysql().Get err: %v", err)
			return nil, err
		}

		var ids []interface{}
		for _, value := range retrieveIDs(category) {
			ids = append(ids, value)
		}
		searchReq.CategoryIDs = ids
	}

	goodsList, err := gs.searchData.Goods().Search(ctx, &searchReq)
	if err != nil {
		log.Errorf("serachdata.NewMysql().Search err: %v", err)
		return nil, err
	}

	log.Debugf("Search es data: %v", goodsList)

	goodsIDs := []uint64{}
	for _, value := range goodsList.Items {
		goodsIDs = append(goodsIDs, uint64(value.ID))
	}

	//通过id批量查询mysql数据
	goods, err := gs.data.NewMysql().Goods().ListByIDs(ctx, goodsIDs, orderby)
	if err != nil {
		log.Errorf("data.NewMysql().ListByIDs err: %v", err)
		return nil, err
	}
	var ret dto.GoodsDTOList
	ret.TotalCount = int(goodsList.TotalCount)
	for _, value := range goods.Items {
		ret.Items = append(ret.Items, &dto.GoodsDTO{
			GoodsDO: *value,
		})
	}
	return &ret, nil
}

func (gs *goodsService) Get(ctx context.Context, ID uint64) (*dto.GoodsDTO, error) {
	goods, err := gs.data.NewMysql().Goods().Get(ctx, ID)
	if err != nil {
		log.Errorf("data.NewMysql().Get err: %v", err)
		return nil, err
	}
	return &dto.GoodsDTO{
		GoodsDO: *goods,
	}, nil
}

// Create TODO canal
func (gs *goodsService) Create(ctx context.Context, goods *v1.GoodsInfo) error {
	/*
		数据先写mysql，然后写es
	*/
	_, err := gs.data.NewMysql().Brands().Get(ctx, uint64(goods.GoodsDO.BrandsID))
	if err != nil {
		return err
	}

	_, err = gs.data.NewMysql().Categorys().Get(ctx, uint64(goods.GoodsDO.CategoryID))
	if err != nil {
		return err
	}

	txn := gs.data.NewMysql().Begin()
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()

	err = gs.data.NewMysql().Goods().CreateInTxn(ctx, txn, goods)
	if err != nil {
		log.Errorf("data.NewMysql().CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}
	model := goods.GoodsDO
	searchDO := do.GoodsSearchDO{
		ID:          model.ID,
		CategoryID:  model.CategoryID,
		BrandsID:    model.BrandsID,
		Name:        model.Name,
		ClickNum:    model.ClickNum,
		FavNum:      model.FavNum,
		MarketPrice: model.MarketPrice,
		GoodsBrief:  model.GoodsBrief,
		ShopPrice:   model.ShopPrice,
	}
	if model.OnSale != nil {
		searchDO.OnSale = *model.OnSale
	}
	if model.ShipFree != nil {
		searchDO.ShipFree = *model.ShipFree
	}
	if model.IsNew != nil {
		searchDO.IsNew = *model.IsNew
	}
	if model.IsHot != nil {
		searchDO.IsHot = *model.IsHot
	}

	// Canal会自动监听binlog，无需在Create方法中主动发送MQ
	// 同步逻辑由Canal监听器（后台goroutine）处理

	return txn.Commit().Error

}

func (gs *goodsService) Update(ctx context.Context, goods *v1.GoodsInfo) error {
	if goods.GoodsDO.BrandsID != 0 {
		_, err := gs.data.NewMysql().Brands().Get(ctx, uint64(goods.GoodsDO.BrandsID))
		if err != nil {
			return err
		}
	}
	if goods.GoodsDO.CategoryID != 0 {
		_, err := gs.data.NewMysql().Categorys().Get(ctx, uint64(goods.GoodsDO.CategoryID))
		if err != nil {
			return err
		}
	}

	txn := gs.data.NewMysql().Begin()
	defer func() { // 很重要
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()

	err := gs.data.NewMysql().Goods().UpdateInTxn(ctx, txn, goods)
	if err != nil {
		return err
	}

	return txn.Commit().Error
}

func (gs *goodsService) CreateInTxn(ctx context.Context, goods *v1.GoodsInfo) error {
	/*
		数据先写mysql，然后写es
	*/
	_, err := gs.data.NewMysql().Brands().Get(ctx, uint64(goods.GoodsDO.BrandsID))
	if err != nil {
		return err
	}

	_, err = gs.data.NewMysql().Categorys().Get(ctx, uint64(goods.GoodsDO.CategoryID))
	if err != nil {
		return err
	}

	// 分布式事务， 异构数据库的事务， 基于可靠消息最终一致性
	// TODO canal
	txn := gs.data.NewMysql().Begin() // 非常小心， 这种方案 也有问题
	defer func() {                    // 很重要
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()

	err = gs.data.NewMysql().Goods().CreateInTxn(ctx, txn, goods)
	if err != nil {
		log.Errorf("data.NewMysql().CreateInTxn err: %v", err)
		txn.Rollback()
		return err
	}
	model := goods.GoodsDO
	searchDO := do.GoodsSearchDO{
		ID:          model.ID,
		CategoryID:  model.CategoryID,
		BrandsID:    model.BrandsID,
		Name:        model.Name,
		ClickNum:    model.ClickNum,
		FavNum:      model.FavNum,
		MarketPrice: model.MarketPrice,
		GoodsBrief:  model.GoodsBrief,
		ShopPrice:   model.ShopPrice,
	}
	if model.OnSale != nil {
		searchDO.OnSale = *model.OnSale
	}
	if model.ShipFree != nil {
		searchDO.ShipFree = *model.ShipFree
	}
	if model.IsNew != nil {
		searchDO.IsNew = *model.IsNew
	}
	if model.IsHot != nil {
		searchDO.IsHot = *model.IsHot
	}

	err = gs.searchData.Goods().Create(ctx, &searchDO) //这个接口如果超时了
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit().Error

}

func (gs *goodsService) UpdateInTxn(ctx context.Context, goods *v1.GoodsInfo) error {
	if goods.GoodsDO.BrandsID != 0 {
		_, err := gs.data.NewMysql().Brands().Get(ctx, uint64(goods.GoodsDO.BrandsID))
		if err != nil {
			return err
		}
	}
	if goods.GoodsDO.CategoryID != 0 {
		_, err := gs.data.NewMysql().Categorys().Get(ctx, uint64(goods.GoodsDO.CategoryID))
		if err != nil {
			return err
		}
	}

	txn := gs.data.NewMysql().Begin() // 非常小心， 这种方案 也有问题
	defer func() {                    // 很重要
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("goodsService.Create panic: %v", err)
			return
		}
	}()

	err := gs.data.NewMysql().Goods().UpdateInTxn(ctx, txn, goods)
	if err != nil {
		return err
	}
	// search 层   (es)
	model := goods.GoodsDO
	searchDO := do.GoodsSearchDO{
		ID:          model.ID,
		CategoryID:  model.CategoryID,
		BrandsID:    model.BrandsID,
		Name:        model.Name,
		ClickNum:    model.ClickNum,
		FavNum:      model.FavNum,
		MarketPrice: model.MarketPrice,
		GoodsBrief:  model.GoodsBrief,
		ShopPrice:   model.ShopPrice,
	}
	if model.OnSale != nil {
		searchDO.OnSale = *model.OnSale
	}
	if model.ShipFree != nil {
		searchDO.ShipFree = *model.ShipFree
	}
	if model.IsNew != nil {
		searchDO.IsNew = *model.IsNew
	}
	if model.IsHot != nil {
		searchDO.IsHot = *model.IsHot
	}

	err = gs.searchData.Goods().Update(ctx, &searchDO) //这个接口如果超时了
	if err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit().Error

}

func (gs *goodsService) Delete(ctx context.Context, ID uint64) error {
	err := gs.data.NewMysql().Goods().Delete(ctx, ID)
	return err
}

func (gs *goodsService) BatchGet(ctx context.Context, ids []uint64) ([]*dto.GoodsDTO, error) {
	// TODO go-zero 非常好用  一次性启动多个goroutine
	var ret []*dto.GoodsDTO
	var callFuncs []func() error
	var mu sync.Mutex
	for _, value := range ids {
		tmp := value
		callFuncs = append(callFuncs, func() error {
			goodsDTO, err := gs.Get(ctx, tmp)
			mu.Lock()
			ret = append(ret, goodsDTO)
			mu.Unlock()
			return err
		})
	}

	err := mr.Finish(callFuncs...)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

var _ GoodsSrv = &goodsService{}
