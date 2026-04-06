package server

import (
	"fmt"
	"maphoto/internal/app/web/api"
	"maphoto/internal/app/web/jwt"
	"maphoto/internal/app/web/view"
	"maphoto/internal/env"
	"maphoto/internal/util"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "maphoto/docs"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	// 初始化本地存储目录
	if err := api.InitStorage(); err != nil {
		util.Logger.Error("初始化存储目录失败: " + err.Error())
	}
}

var (
	Config = jwt.Config{
		JWTScrect: util.ShortUID(12),
	}
)

func MakeRouter(env *env.Options) *echo.Echo {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://www.wsh233.cn", "https://wsh233.cn"},
		AllowCredentials: true,
		AllowMethods: []string{
			echo.GET,
			echo.HEAD,
			echo.PUT,
			echo.PATCH,
			echo.POST,
			echo.DELETE,
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	TokenLookup: "header:X-XSRF-TOKEN",
	// }))

	// 全局中间件：为所有路由设置配置项
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", Config)
			return next(c)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return e
}

func RunServer(env *env.Options) {
	jwt.InitJWT(env)
	e := MakeRouter(env)
	publicGroup := e.Group(env.UrlPrefix)
	protectGroup := e.Group(env.UrlPrefix)
	protectGroup.Use(jwt.DevJWTInjector)
	protectGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(Config.JWTScrect),
		Skipper: func(c echo.Context) bool {
			cPath := c.Path()
			switch cPath {
			case "/api/v1/user/login":
				return true
			default:
				return false
			}
		},
	}))

	if env.Debug {
		publicGroup.GET("/swagger/*", echoSwagger.WrapHandler)
		fmt.Printf("swagger at:http://127.0.0.1:%d%s/swagger/index.html\n", env.Port, env.UrlPrefix)

	}

	api.InitAPI(protectGroup, publicGroup, env)
	view.InitPage(publicGroup)

	// 静态文件服务：访问上传的图片
	e.Static("/uploads", api.UploadDir)
	url := fmt.Sprintf("http://127.0.0.1:%d%s", env.Port, env.UrlPrefix)
	util.OpenBrowser(fmt.Sprintf("%s/#/admin", url))
	fmt.Printf("app server at:%s\n", url)
	e.Logger.Fatal(e.Start(fmt.Sprintf("127.0.0.1:%d", env.Port)))

}
