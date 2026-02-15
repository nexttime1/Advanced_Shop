package admin

import (
	"Advanced_Shop/app/xshop/api/config"
	"Advanced_Shop/app/xshop/api/internal/controller/goods/v1"
	v13 "Advanced_Shop/app/xshop/api/internal/controller/order/v1
	v12 "Advanced_Shop/app/xshop/api/internal/controller/sms/v1"
	"Advanced_Shop/app/xshop/api/internal/controller/user/v1"
	"Advanced_Shop/app/xshop/api/internal/data/rpc"
	"Advanced_Shop/app/xshop/api/internal/service"
	"Advanced_Shop/gnova/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	// 服务
	userGroup := g.Group("/u")
	// 版本
	v1 := userGroup.Group("/v1")
	jwtAuth := newJWTAuth(cfg.Jwt)
	// 路由前缀
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

	// 服务
	goodGroup := g.Group("/g")
	// 版本
	v1 = goodGroup.Group("/v1")
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

	// 订单路由
	// 服务
	orderGroup := g.Group("/o")
	// 版本
	v1 = orderGroup.Group("/v1")
	orderRouter := v1.Group("orders")
	{
		orderController := v13.NewOrderController(serviceFactory, g.Translator(), cfg.Aliyun)

		{
			// order 相关
			orderRouter.GET("", jwtAuth.AuthFunc(), orderController.OrderListView)       // 查看所有订单
			orderRouter.POST("", jwtAuth.AuthFunc(), orderController.OrderCreateView)    // 创建订单
			orderRouter.GET("/:id", jwtAuth.AuthFunc(), orderController.OrderDetailView) // 订单细节
		}
		// cart 相关
		cartRouter := v1.Group("shopcarts")
		{
			cartRouter.GET("", orderController.CartListView)              // 购物车列表
			cartRouter.DELETE("/:id", orderController.DeleteCartItemView) // 删除条目
			cartRouter.POST("", orderController.AddItemView)              // 添加商品到购物车
			cartRouter.PATCH("/:id", orderController.UpdatePatchView)     // 更新购物车中的某个商品
		}
		// 阿里云回调
		alipayRouter := v1.Group("/pay")
		alipayRouter.POST("/callback", orderController.AlipayCallBackView)

	}

}
