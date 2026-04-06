package view

import "github.com/labstack/echo/v4"

func InitPage(gp *echo.Group) {
	gp.GET("/", PageIndex)
	gp.GET("/static/:filename", AssetsFinder)
	gp.GET("/photo/:name", ViewPhoto)
	gp.GET("/thumbnail/:name", ViewPhoto)
}
