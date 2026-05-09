package admin

import (
	"Advanced_Shop/app/pkg/common"
	"Advanced_Shop/app/xshop/api/config"
	v2 "Advanced_Shop/app/xshop/api/internal/controller/action/v1"
	"Advanced_Shop/app/xshop/api/internal/controller/goods/v1"
	v3 "Advanced_Shop/app/xshop/api/internal/controller/order/v1"
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
		ugroup.POST("login", uController.Login)
		ugroup.POST("register", common.Wrapper(uController.Register))

		ugroup.GET("detail", jwtAuth.AuthFunc(), common.Wrapper(uController.GetUserDetail))
		ugroup.GET("list", jwtAuth.AuthFunc(), common.Wrapper(uController.UserListView))
		ugroup.PATCH("update", jwtAuth.AuthFunc(), common.Wrapper(uController.GetUserDetail))
	}

	baseRouter := v1.Group("base")
	{
		smsCtl := v12.NewSmsController(serviceFactory, g.Translator())
		baseRouter.POST("send_sms", common.Wrapper(smsCtl.SendSms))
		baseRouter.GET("captcha", user.GetCaptcha)
	}

	//商品相关的api

	// 服务
	goodGroup := g.Group("/g")
	// 版本
	v1 = goodGroup.Group("/v1")
	goodsRouter := v1.Group("good")
	{
		goodsController := goods.NewGoodsController(serviceFactory, g.Translator())
		// 商品相关
		goodsRouter.GET("/list", common.Wrapper(goodsController.GetGoodListView)) // 限流
		goodsRouter.POST("/", common.Wrapper(goodsController.CreateGoodView))
		goodsRouter.GET("/:id", common.Wrapper(goodsController.GoodDetailView))
		goodsRouter.PUT("/:id", common.Wrapper(goodsController.GoodUpdateView))
		goodsRouter.PATCH("/:id", common.Wrapper(goodsController.GoodPatchUpdateView))
		goodsRouter.DELETE("/:id", common.Wrapper(goodsController.GoodDeleteView))

		// 图片相关
		v1.GET("banners", common.Wrapper(goodsController.GetBannerListView))
		v1.POST("banners", common.Wrapper(goodsController.CreateBannerView))
		v1.PUT("banners/:id", common.Wrapper(goodsController.UpdateBannerView))
		v1.DELETE("banners/:id", common.Wrapper(goodsController.DeleteBannerView))

		// 分类相关
		v1.GET("categorys", common.Wrapper(goodsController.GetAllCategoryView))
		v1.GET("categorys/:id", common.Wrapper(goodsController.GetSubCategoryView))
		v1.POST("categorys", common.Wrapper(goodsController.CreateCategoryView))
		v1.PUT("categorys/:id", common.Wrapper(goodsController.UpdateCategoryView))
		v1.DELETE("categorys/:id", common.Wrapper(goodsController.DeleteCategoryView))

		// 品牌相关
		v1.GET("brands", common.Wrapper(goodsController.BrandListView))
		v1.POST("brands", common.Wrapper(goodsController.CreateBrandView))
		v1.PUT("brands/:id", common.Wrapper(goodsController.UpdateBrandView))
		v1.DELETE("brands/:id", common.Wrapper(goodsController.DeleteBrandView))

		// 第三张表
		v1.GET("categorybrands", common.Wrapper(goodsController.CategoryBrandListView))    //所有的 第三张表
		v1.GET("categorybrands/:id", common.Wrapper(goodsController.CategoryAllBrandView)) //某个分类下的所有品牌
		v1.POST("categorybrands", common.Wrapper(goodsController.CreateCategoryBrandView))
		v1.PUT("categorybrands/:id", common.Wrapper(goodsController.UpdateCategoryBrandView))
		v1.DELETE("categorybrands/:id", common.Wrapper(goodsController.DeleteCategoryBrandView))
	}

	// 订单路由
	// 服务
	orderGroup := g.Group("/o")
	// 版本
	v1 = orderGroup.Group("/v1")
	orderRouter := v1.Group("orders")
	{
		orderController := v3.NewOrderController(serviceFactory, g.Translator(), cfg.Aliyun)

		{
			// order 相关
			orderRouter.GET("", jwtAuth.AuthFunc(), common.Wrapper(orderController.OrderListView))       // 查看所有订单
			orderRouter.POST("", jwtAuth.AuthFunc(), common.Wrapper(orderController.OrderCreateView))    // 创建订单
			orderRouter.GET("/:id", jwtAuth.AuthFunc(), common.Wrapper(orderController.OrderDetailView)) // 订单细节
		}
		// cart 相关
		cartRouter := v1.Group("shopcarts")
		{
			cartRouter.GET("", common.Wrapper(orderController.CartListView))              // 购物车列表
			cartRouter.DELETE("/:id", common.Wrapper(orderController.DeleteCartItemView)) // 删除条目
			cartRouter.POST("", common.Wrapper(orderController.AddItemView))              // 添加商品到购物车
			cartRouter.PATCH("/:id", common.Wrapper(orderController.UpdatePatchView))     // 更新购物车中的某个商品
		}
		// 阿里云回调
		alipayRouter := v1.Group("/pay")
		alipayRouter.POST("/callback", orderController.AlipayCallBackView)

	}

	// 服务
	addressGroup := g.Group("/up")
	v1 = addressGroup.Group("/v1")
	// 地址核心路由组
	addressRouter := v1.Group("address")
	// 实例化地址控制器（请根据你的实际包路径调整）
	ActionController := v2.NewActionController(serviceFactory, g.Translator())
	{

		// 地址相关接口（添加认证和Trace中间件）
		addressRouter.GET("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.AddressListView))          // 查看所有地址
		addressRouter.DELETE("/:id", jwtAuth.AuthFunc(), common.Wrapper(ActionController.DeleteAddressView)) // 删除地址
		addressRouter.POST("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.AddressCreateView))       // 创建地址
		addressRouter.PUT("/:id", jwtAuth.AuthFunc(), common.Wrapper(ActionController.UpdateAddressView))    // 修改地址
	}

	// 收藏模块路由
	// 服务分组
	collectionRouter := v1.Group("userfavs")
	{

		collectionRouter.GET("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.CollectionListView))               // 查看收藏
		collectionRouter.DELETE("/:good_id", jwtAuth.AuthFunc(), common.Wrapper(ActionController.CollectionDeleteView)) // 删除收藏
		collectionRouter.POST("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.CollectionAddView))               // 添加收藏
		collectionRouter.GET("/:good_id", jwtAuth.AuthFunc(), common.Wrapper(ActionController.CollectionDetailView))    // 查看详情
	}

	// 消息核心路由组
	messageRouter := v1.Group("message")
	{

		// 消息相关接口（添加认证和Trace中间件）
		messageRouter.GET("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.MessageListView))    // 消息列表
		messageRouter.POST("", jwtAuth.AuthFunc(), common.Wrapper(ActionController.CreateMessageView)) // 添加留言
	}

}
