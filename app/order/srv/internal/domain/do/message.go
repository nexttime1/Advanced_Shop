package do

type MQMessageType int

const (
	DirectPass MQMessageType = 0
	OptionFail MQMessageType = 1
	Continuing MQMessageType = 2
)
