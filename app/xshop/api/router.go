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
		// 商品相关
		goodsRouter.GET("good/list", goodsController.GetGoodListView) // 限流
		goodsRouter.POST("good", goodsController.CreateGoodView)
		goodsRouter.GET("good/:id", goodsController.GoodDetailView)
		goodsRouter.PUT("good/:id", goodsController.GoodUpdateView)
		goodsRouter.PATCH("good/:id", goodsController.GoodPatchUpdateView)
		goodsRouter.DELETE("good/:id", goodsController.GoodDeleteView)

		// 图片相关
		goodsRouter.GET("banners", goodsController.GetBannerListView)
		goodsRouter.POST("banners", goodsController.CreateBannerView)
		goodsRouter.PUT("banners/:id", goodsController.UpdateBannerView)
		goodsRouter.DELETE("banners/:id", goodsController.DeleteBannerView)

		// 分类相关
		goodsRouter.GET("categorys", goodsController.GetAllCategoryView)
		goodsRouter.GET("categorys/:id", goodsController.GetSubCategoryView)
		goodsRouter.POST("categorys", goodsController.CreateCategoryView)
		goodsRouter.PUT("categorys/:id", goodsController.UpdateCategoryView)
		goodsRouter.DELETE("categorys/:id", goodsController.DeleteCategoryView)

		// 品牌相关
		goodsRouter.GET("brands", goodsController.BrandListView)
		goodsRouter.POST("brands", goodsController.CreateBrandView)
		goodsRouter.PUT("brands/:id", goodsController.UpdateBrandView)
		goodsRouter.DELETE("brands/:id", goodsController.DeleteBrandView)

		// 第三张表
		goodsRouter.GET("categorybrands", goodsController.CategoryBrandListView)    //所有的 第三张表
		goodsRouter.GET("categorybrands/:id", goodsController.CategoryAllBrandView) //某个分类下的所有品牌
		goodsRouter.POST("categorybrands", goodsController.CreateCategoryBrandView)
		goodsRouter.PUT("categorybrands/:id", goodsController.UpdateCategoryBrandView)
		goodsRouter.DELETE("categorybrands/:id", goodsController.DeleteCategoryBrandView)
	}
}
