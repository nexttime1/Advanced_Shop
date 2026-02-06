package admin

import (
	"Advanced_Shop/app/xshop/api/config"
	"Advanced_Shop/app/xshop/api/internal/controller/goods/v1"
	v12 "Advanced_Shop/app/xshop/api/internal/controller/sms/v1"
	"Advanced_Shop/app/xshop/api/internal/controller/user/v1"
	"Advanced_Shop/app/xshop/api/internal/data/rpc"
	"Advanced_Shop/app/xshop/api/internal/service"
	"Advanced_Shop/gnova/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	ugroup := v1.Group("/user")

	data, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(err)
	}

	serviceFactory := service.NewService(data, cfg.Sms, cfg.Jwt)
	uController := user.NewUserController(g.Translator(), serviceFactory)
	{
		ugroup.POST("pwd_login", uController.Login)
		ugroup.POST("register", uController.Register)

		jwtAuth := newJWTAuth(cfg.Jwt)
		ugroup.GET("detail", jwtAuth.AuthFunc(), uController.GetUserDetail)
		ugroup.PATCH("update", jwtAuth.AuthFunc(), uController.GetUserDetail)
	}

	baseRouter := v1.Group("base")
	{
		smsCtl := v12.NewSmsController(serviceFactory, g.Translator())
		baseRouter.POST("send_sms", smsCtl.SendSms)
		baseRouter.GET("captcha", user.GetCaptcha)
	}

	//商品相关的api
	goodsRouter := v1.Group("goods")
	{
		goodsController := goods.NewGoodsController(serviceFactory, g.Translator())
		goodsRouter.GET("", goodsController.List)
	}
}
