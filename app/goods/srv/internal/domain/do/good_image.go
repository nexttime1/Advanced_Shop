package do

import (
	gorm2 "Advanced_Shop/app/pkg/gorm"
	"Advanced_Shop/pkg/errors"
	"database/sql/driver"
	"go.uber.org/zap"
)

type ImageType int

const (
	MainImageType   ImageType = 1
	DetailImageType ImageType = 2
	OtherImageType  ImageType = 3
)

func (t ImageType) Value() (driver.Value, error) {
	return int64(t), nil // 改为返回int64，符合driver.Value规范
}

// Scan 实现 sql.Scanner 接口：将数据库类型（tinyint）转 Go 类型
func (t *ImageType) Scan(value interface{}) error {
	// 兼容更多类型（比如int/uint），避免转换失败
	var val int64
	switch v := value.(type) {
	case int64:
		val = v
	case int:
		val = int64(v)
	case uint:
		val = int64(v)
	case uint64:
		val = int64(v)
	default:
		zap.S().Errorf("图片类型转换失败，不支持的类型: %T, 值: %v", value, value)
		return errors.New("invalid image type value")
	}
	*t = ImageType(val)
	return nil
}

// IsValid 校验图片类型是否合法（oneof 1/2/3）
func (t ImageType) IsValid() bool {
	return t == MainImageType || t == DetailImageType || t == OtherImageType
}

// String 转字符串描述（方便日志打印/返回给前端）
func (t ImageType) String() string {
	switch t {
	case MainImageType:
		return "main"
	case DetailImageType:
		return "detail"
	case OtherImageType:
		return "other"
	default:
		return "unknown"
	}
}

type GoodsImageModel struct {
	gorm2.Model
	GoodsID   int32     `gorm:"type:int;not null;comment:商品ID（逻辑外键，关联good_models.id）;index:idx_goods_image_goods"`
	ImageURL  string    `gorm:"type:varchar(255);not null;comment:图片访问URL（七牛云）"`
	Sort      int       `gorm:"type:int;not null;default:0;comment:排序序号（越小越靠前）"`
	IsMain    bool      `gorm:"default:false;not null;comment:是否主图（一个商品仅一个主图）"`
	ImageType ImageType `gorm:"type:tinyint(1);not null;default:3;comment:图片类型（1=主图，2=详情图，3=其他）"`
}

func (GoodsImageModel) TableName() string {
	return "goods_image_models"
}
