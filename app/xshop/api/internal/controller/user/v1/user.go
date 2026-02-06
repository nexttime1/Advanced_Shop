package user

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop/app/mxshop/api/internal/service"
)

type userServer struct {
	trans ut.Translator

	sf service.ServiceFactory
}

func NewUserController(trans ut.Translator, sf service.ServiceFactory) *userServer {
	return &userServer{trans, sf}
}
