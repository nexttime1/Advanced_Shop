package user

import (
	"Advanced_Shop/app/xshop/api/internal/service"
	ut "github.com/go-playground/universal-translator"
)

type userServer struct {
	trans ut.Translator

	sf service.ServiceFactory
}

func NewUserController(trans ut.Translator, sf service.ServiceFactory) *userServer {
	return &userServer{trans, sf}
}
