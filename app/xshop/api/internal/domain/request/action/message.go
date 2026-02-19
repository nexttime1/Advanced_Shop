package action

type MessageRequest struct {
	MessageType int32  `form:"type" json:"type" binding:"required,oneof=1 2 3 4 5"`
	Subject     string `form:"subject" json:"subject" binding:"required"`
	Message     string `form:"message" json:"message" binding:"required"`
	File        string `form:"file" json:"file" binding:"required"`
}

type MessageResponse struct {
	Id          int32  `json:"id"`
	UserId      int32  `json:"user_id"`
	MessageType int32  `json:"message_type"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	File        string `json:"file"`
}
