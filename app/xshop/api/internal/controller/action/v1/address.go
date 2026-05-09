package v1

import (
	proto "Advanced_Shop/api/action/v1"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/app/xshop/api/internal/domain/request/action"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

func (ac *actionController) AddressListView(c *gin.Context) error {
	log.Info("address list function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	ctx := c.Request.Context()
	list, err := ac.srv.Address().GetAddressList(ctx, &proto.AddressRequest{
		UserId: userID,
	})
	if err != nil {
		return err
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
	common.OkWithList(c, response, list.Total)
	return nil
}

func (ac *actionController) AddressCreateView(c *gin.Context) error {
	log.Info("address create function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr action.AddressCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

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
		return err
	}
	response := action.AddressCreateResponse{
		Id: address.Id,
	}

	common.OkWithData(c, response)
	return nil
}

func (ac *actionController) DeleteAddressView(c *gin.Context) error {
	log.Info("address delete function called ...")
	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var cr action.AddressIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)
	}

	ctx := c.Request.Context()
	_, err = ac.srv.Address().DeleteAddress(ctx, &proto.AddressRequest{
		Id:     cr.Id,
		UserId: userID,
	})
	if err != nil {
		return err
	}
	common.OkWithMessage(c, "删除成功")
	return nil
}

func (ac *actionController) UpdateAddressView(c *gin.Context) error {
	log.Info("address update function called ...")

	userID, _, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}

	var idRequest action.AddressIdRequest
	err = c.ShouldBindUri(&idRequest)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

	}

	var cr action.AddressUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		return gin2.HandleValidatorError(c, err, ac.trans)

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
		return err

	}
	common.OkWithMessage(c, "更新成功")
	return nil
}
