package db

import (
	v1 "Advanced_Shop/app/action/srv/internal/data/v1"
	"Advanced_Shop/app/action/srv/internal/domain/do"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"gorm.io/gorm"
)

type messageData struct {
	db *gorm.DB
}

func newMessage(factory *mysqlFactory) v1.MessageStore {
	return &messageData{
		db: factory.db,
	}
}

// ListByUserID 根据用户ID获取留言列表
func (s *messageData) ListByUserID(ctx context.Context, userID int32) ([]*do.LeavingMessageDO, int64, error) {
	var messages []*do.LeavingMessageDO
	tx := s.db.WithContext(ctx).Model(&do.LeavingMessageDO{}).Where("user_id = ?", userID)

	// 统计总数
	var count int64
	if err := tx.Count(&count).Error; err != nil {
		log.Errorf("Message ListByUserID count err:%v", err)
		return nil, 0, errors.WithCode(code.ErrMessageQuery, err.Error())
	}

	// 查询列表
	if err := tx.Find(&messages).Error; err != nil {
		log.Errorf("Message ListByUserID find err:%v", err)
		return nil, 0, errors.WithCode(code.ErrMessageQuery, err.Error())
	}

	return messages, count, nil
}

// Create 创建留言
func (s *messageData) Create(ctx context.Context, message *do.LeavingMessageDO) error {
	err := s.db.WithContext(ctx).Create(message).Error
	if err != nil {
		log.Errorf("Message Create err:%v", err)
		return errors.WithCode(code.ErrMessageCreate, err.Error())
	}
	return nil
}

// 确保实现了接口
var _ v1.MessageStore = &messageData{}
