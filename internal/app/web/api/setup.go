package api

import (
	"maphoto/internal/env"

	"github.com/labstack/echo/v4"
)

var (
	Env *env.Options
)

// 集中绑定接口
func InitAPI(protectGroup *echo.Group, publickGroup *echo.Group, env *env.Options) {
	Env = env
	protectGroup.GET("/api/v1/place/all", PlaceAll)
	protectGroup.GET("/api/v1/user/config", UserConfig)
	protectGroup.GET("/api/v1/user/logout", UserLogout)
	protectGroup.POST("/api/v1/user/config/update", UserConfigGetUpate)
	protectGroup.POST("/api/v1/place/add", PlaceAdd)
	protectGroup.POST("/api/v1/place/update/:id", PlaceUpdate)
	protectGroup.POST("/api/v1/place/cover", PlaceCover)
	protectGroup.POST("/api/v1/place/pic/add", PlacePicAdd)
	protectGroup.POST("/api/v1/place/pic/del", PlacePicDel)
	protectGroup.GET("/api/v1/place/del/:id", PlaceDel)
	protectGroup.GET("/api/v1/place/init/:username", PlaceInit)

	// User management APIs
	protectGroup.GET("/api/v1/users", UserList)
	protectGroup.POST("/api/v1/user/add", UserAdd)
	protectGroup.POST("/api/v1/user/update/:id", UserUpdate)
	protectGroup.GET("/api/v1/user/del/:id", UserDelete)

	publickGroup.POST("/api/v1/user/login", UserLogin)
	publickGroup.GET("/api/v1/share/:username", UserShare)

	// Storage management APIs
	protectGroup.POST("/api/v1/storage/upload", StorageUpload)
	protectGroup.GET("/api/v1/storage/files", StorageList)
	protectGroup.DELETE("/api/v1/storage/files/:folder/:filename", StorageDelete)
	protectGroup.POST("/api/v1/storage/files/:folder/:filename/rename", StorageRename)
}
