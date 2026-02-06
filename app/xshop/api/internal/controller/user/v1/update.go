package user

import (
	"time"

	"github.com/gin-gonic/gin"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/pkg/common/core"
	jtime "mxshop/pkg/common/time"
)

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}

func (us *userServer) UpdateUser(ctx *gin.Context) {
	updateForm := UpdateUserForm{}
	if err := ctx.ShouldBind(&updateForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	userID, _ := ctx.Get(middlewares.KeyUserID)
	userIDInt := uint64(userID.(float64))
	userDTO, err := us.sf.Users().Get(ctx, userIDInt)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	userDTO.NickName = updateForm.Name

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateForm.Birthday, loc)
	userDTO.NickName = updateForm.Name
	userDTO.Birthday = jtime.Time{birthDay}
	userDTO.Gender = updateForm.Gender
	err = us.sf.Users().Update(ctx, userDTO)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)
}
