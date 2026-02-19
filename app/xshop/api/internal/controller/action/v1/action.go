package v1

import (
	"Advanced_Shop/app/xshop/api/internal/service"
	ut "github.com/go-playground/universal-translator"
)

type actionController struct {
	trans ut.Translator
	srv   service.ServiceFactory
}

func NewActionController(srv service.ServiceFactory, trans ut.Translator) *actionController {
	return &actionController{
		srv:   srv,
		trans: trans,
	}
}
