package user

import (
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/common"
	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	Id       int32  `json:"id"`
	Password string `json:"password"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	BirthDay string `json:"birth_day"`
	Gender   string `json:"gender"`
	Role     int    `json:"role"`
}

func (us *userServer) UserListView(c *gin.Context) error {
	log.Info("UserListView is called")
	_, role, err := common.GetAuthUser(c)
	if err != nil {
		return err
	}
	if role != 1 {
		return errors.WithCode(code.ErrInsufficientPermissions, "权限不足")
	}

	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		return gin2.HandleValidatorError(c, err, us.trans)

	}
	ctx := c.Request.Context()
	userListResponse, err := us.sf.Users().GetList(ctx, common.PageInfo{
		Limit: cr.Limit,
		Page:  cr.Page,
	})
	if err != nil {
		return err
	}
	var response []UserListResponse
	for _, v := range userListResponse.Items {
		password := v.PassWord[0:1] + "*****"
		birthDayStr := v.Birthday.Format("2006-01-02")
		response = append(response, UserListResponse{
			Id:       int32(v.ID),
			Password: password,
			Mobile:   v.Mobile,
			NickName: v.NickName,
			BirthDay: birthDayStr,
			Gender:   v.Gender,
			Role:     int(v.Role),
		})
	}
	common.OkWithList(c, response, int32(userListResponse.TotalCount))
	return nil
}
