package core

import (
	"fmt"
	"net/http"

	"Advanced_Shop/pkg/errors"
	"github.com/gin-gonic/gin"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"msg"`

	Detail string `json:"detail"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteErrResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteErrResponse(c *gin.Context, err error, data interface{}) {
	errStr := fmt.Sprintf("%#+v", err)
	coder := errors.ParseCoder(err)
	c.JSON(coder.HTTPStatus(), ErrResponse{
		Code:      coder.Code(),
		Message:   coder.String(),
		Detail:    errStr,
		Reference: coder.Reference(),
	})

	return

}

type SuccessBaseResponse struct {
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type DataListResponse struct {
	List  any   `json:"list"`
	Count int32 `json:"count"`
}

var empty = map[string]interface{}{}

func response(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, SuccessBaseResponse{
		Data: data,
		Msg:  msg,
	})
}

func OK(c *gin.Context, data interface{}, msg string) {
	response(c, data, msg)
}

func OkWithMessage(c *gin.Context, msg string) {
	response(c, empty, msg)
}

func OkWithData(c *gin.Context, data interface{}) {
	response(c, data, "成功")
}

func OkWithList(c *gin.Context, list interface{}, count int32) {
	response(c, DataListResponse{
		List:  list,
		Count: count,
	}, "成功")

}
