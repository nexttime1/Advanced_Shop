package v1

import (
	proto "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/action"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ac *actionController) MessageListView(c *gin.Context) error {
	log.Info("message list function called.")
	userID, role, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	request := proto.MessageRequest{
		UserId: userID,
	}
	if role == 2 {
		request.UserId = 0
	}

	ctx := c.Request.Context()
	List, err := ac.srv.Message().MessageList(ctx, &request)

	if err != nil {
		return err

	}

	var response []action.MessageResponse

	for _, model := range List.Data {
		response = append(response, action.MessageResponse{
			Id:          model.Id,
			UserId:      model.UserId,
			MessageType: model.MessageType,
			Subject:     model.Subject,
			Message:     model.Message,
			File:        model.File,
		})

	}

	common.OkWithList(c, response, List.Total)
	return nil
}

func (ac *actionController) CreateMessageView(c *gin.Context) error {
	log.Info("message create function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr action.MessageRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

	}

	ctx := c.Request.Context()
	req, err := ac.srv.Message().CreateMessage(ctx, &proto.MessageRequest{
		UserId:      userID,
		MessageType: cr.MessageType,
		Subject:     cr.Subject,
		Message:     cr.Message,
		File:        cr.File,
	})
	if err != nil {
		return err
	}
	RMap := map[string]interface{}{
		"id": req.Id,
	}
	common.OkWithData(c, RMap)
	return nil
}
