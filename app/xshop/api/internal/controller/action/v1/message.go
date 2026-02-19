package v1

import (
	proto "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/action"
	"Advanced_Shop/pkg/common/core"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ac *actionController) MessageListView(c *gin.Context) {
	log.Info("message list function called.")
	userID, role, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
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
		core.WriteErrResponse(c, err, nil)
		return
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

	core.OkWithList(c, response, List.Total)

}

func (ac *actionController) CreateMessageView(c *gin.Context) {
	log.Info("message create function called.")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr action.MessageRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
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
		core.WriteErrResponse(c, err, nil)
		return
	}
	RMap := map[string]interface{}{
		"id": req.Id,
	}
	core.OkWithData(c, RMap)

}
