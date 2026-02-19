package do

import "Advanced_Shop/app/pkg/gorm"

type MessageType int8

const (
	LeaveWord MessageType = 1
	Complaint MessageType = 2
	inquiry   MessageType = 3
	AfterSale MessageType = 4
	AskBuy    MessageType = 5
)

type LeavingMessageDO struct {
	gorm.Model
	UserId      int32       `gorm:"type:int(11);index"`
	MessageType MessageType `gorm:"type:int(11)"`
	Subject     string      `gorm:"type:varchar(128)"`
	Message     string      `gorm:"type:varchar(128)"`
	File        string      `gorm:"type:varchar(200)"`
}

func (LeavingMessageDO) TableName() string {
	return "leaving_message_models"
}
