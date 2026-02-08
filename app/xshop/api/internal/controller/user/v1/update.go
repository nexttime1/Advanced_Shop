package user

import (
	"time"

	gin2 "Advanced_Shop/app/pkg/translator/gin"
	"Advanced_Shop/gnova/server/restserver/middlewares"
	"Advanced_Shop/pkg/common/core"
	jtime "Advanced_Shop/pkg/common/time"
	"github.com/gin-gonic/gin"
)

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
	Password string `json:"password,omitempty"`
}

func (us *userServer) UpdateUser(ctx *gin.Context) {
	var cr UpdateUserForm
	if err := ctx.ShouldBind(&cr); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	userID, _ := ctx.Get(middlewares.KeyUserID)
	userIDInt := uint64(userID.(float64))
	userDTO, err := us.sf.Users().Get(ctx, userIDInt)
	if err != nil {
		core.WriteErrResponse(ctx, err, nil)
		return
	}

	userDTO.NickName = cr.Name
	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", cr.Birthday, loc)
	userDTO.Birthday = jtime.Time{birthDay}
	userDTO.Gender = cr.Gender
	userDTO.PassWord = cr.Password

	err = us.sf.Users().Update(ctx, userDTO)
	if err != nil {
		core.WriteErrResponse(ctx, err, nil)
		return
	}
	core.OkWithMessage(ctx, "修改成功")
}
