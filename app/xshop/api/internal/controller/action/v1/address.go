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

func (ac *actionController) AddressListView(c *gin.Context) {
	log.Info("address list function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}
	ctx := c.Request.Context()
	list, err := ac.srv.Address().GetAddressList(ctx, &proto.AddressRequest{
		UserId: userID,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	var response []action.AddressListResponse
	for _, model := range list.Data {
		response = append(response, action.AddressListResponse{
			Id:           model.Id,
			UserId:       model.UserId,
			Province:     model.Province,
			City:         model.City,
			District:     model.District,
			Address:      model.Address,
			SignerName:   model.SignerName,
			SignerMobile: model.SignerMobile,
		})
	}
	core.OkWithList(c, response, list.Total)

}

func (ac *actionController) AddressCreateView(c *gin.Context) {
	log.Info("address create function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr action.AddressCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}
	ctx := c.Request.Context()
	address, err := ac.srv.Address().CreateAddress(ctx, &proto.AddressRequest{
		UserId:       userID,
		Province:     cr.Province,
		City:         cr.City,
		District:     cr.District,
		Address:      cr.Address,
		SignerName:   cr.SignerName,
		SignerMobile: cr.SignerMobile,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	response := action.AddressCreateResponse{
		Id: address.Id,
	}

	core.OkWithData(c, response)

}

func (ac *actionController) DeleteAddressView(c *gin.Context) {
	log.Info("address delete function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var cr action.AddressIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}

	ctx := c.Request.Context()
	_, err = ac.srv.Address().DeleteAddress(ctx, &proto.AddressRequest{
		Id:     cr.Id,
		UserId: userID,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "删除成功")

}

func (ac *actionController) UpdateAddressView(c *gin.Context) {
	log.Info("address update function called ...")

	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		core.WriteErrResponse(c, errors.New("未登录"), nil)
		return
	}

	var idRequest action.AddressIdRequest
	err = c.ShouldBindUri(&idRequest)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}

	var cr action.AddressUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		gin2.HandleValidatorError(c, err, ac.trans)
		return
	}
	ctx := c.Request.Context()
	_, err = ac.srv.Address().UpdateAddress(ctx, &proto.AddressRequest{
		Id:           idRequest.Id,
		UserId:       userID,
		Province:     cr.Province,
		City:         cr.City,
		District:     cr.District,
		Address:      cr.Address,
		SignerName:   cr.SignerName,
		SignerMobile: cr.SignerMobile,
	})
	if err != nil {
		core.WriteErrResponse(c, err, nil)
		return
	}
	core.OkWithMessage(c, "更新成功")

}
