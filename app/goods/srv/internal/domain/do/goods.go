package do

import (
	"encoding/json"

	gorm2 "Advanced_Shop/app/pkg/gorm"
	"database/sql/driver"
)

type GoodsSearchDO struct {
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	BrandsID   int32 `json:"brands_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

func (GoodsSearchDO) GetIndexName() string {
	return "goods"
}

type GoodsSearchDOList struct {
	TotalCount int64            `json:"totalCount,omitempty"`
	Items      []*GoodsSearchDO `json:"items"`
}

type GoodsDO struct {
	gorm2.Model

	CategoryID int32       `gorm:"type:int;not null;comment:分类ID（逻辑外键）;index:idx_goods_category"`
	Category   *CategoryDO `gorm:"foreignKey:CategoryID;references:ID;constraint:<-:false,foreignKey:no action"`

	BrandsID int32     `gorm:"type:int;not null;comment:品牌ID（逻辑外键）;index:idx_goods_brand"`
	Brands   *BrandsDO `gorm:"foreignKey:BrandsID;references:ID;constraint:<-:false,foreignKey:no action"`

	OnSale      *bool   `gorm:"default:false;not null;comment:是否上架"`
	ShipFree    *bool   `gorm:"default:false;not null;comment:是否包邮"`
	IsNew       *bool   `gorm:"default:false;not null;comment:是否新品"`
	IsHot       *bool   `gorm:"default:false;not null;comment:是否热销"`
	Name        string  `gorm:"type:varchar(50);not null;comment:商品名称;index:idx_goods_name"`
	GoodsSn     string  `gorm:"type:varchar(50);not null;comment:商品编号;uniqueIndex:idx_goods_sn"`
	ClickNum    int32   `gorm:"type:int;default:0;not null;comment:点击量"`
	SoldNum     int32   `gorm:"type:int;default:0;not null;comment:销量"`
	FavNum      int32   `gorm:"type:int;default:0;not null;comment:收藏量"`
	MarketPrice float32 `gorm:"not null;comment:市场价"`
	ShopPrice   float32 `gorm:"not null;comment:售价;index:idx_goods_price"`
	GoodsBrief  string  `gorm:"type:varchar(100);not null;comment:商品简介"`

	// 方便查询商品的所有图片（Gorm虚拟字段，不存数据库）
	Images []*GoodsImageModel `gorm:"foreignKey:GoodsID;references:ID;constraint:<-:false,foreignKey:no action"`
}

func (GoodsDO) TableName() string {
	return "good_models"
}

// GormList 去掉gorm的依赖
type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type GoodsDOList struct {
	TotalCount int64      `json:"totalCount,omitempty"`
	Items      []*GoodsDO `json:"items"`
}
