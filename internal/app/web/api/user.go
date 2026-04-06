package api

import (
	"maphoto/internal/app/form"
	"maphoto/internal/app/model"
	"maphoto/internal/app/store"
	"maphoto/internal/util"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/project"

	jwter "maphoto/internal/app/web/jwt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
)

// @Summary  get all users
// @Tags user
// @Security     JWT
// @Success 200 {object} Response
// @Router /api/v1/users [GET]
func UserList(c echo.Context) error {
	var users []model.User
	result := store.DB.Select("id", "username", "admin", "last_login", "created_at", "updated_at").Find(&users)
	if result.Error != nil {
		return ApiFailed(c, 500, result.Error.Error())
	}
	return ApiSuccess(c, users)
}

// @Summary  create user
// @Tags user
// @Security     JWT
// @Accept       json
// @Param        data  body  form.UserCreateForm true  "data"
// @Success 200 {object} Response
// @Router /api/v1/user/add [POST]
func UserAdd(c echo.Context) error {
	u := new(form.UserCreateForm)
	if err := c.Bind(u); err != nil {
		return ApiFailed(c, 400, err.Error())
	}

	// Check if username already exists
	var existingUser model.User
	store.DB.Find(&existingUser, "username = ?", u.Username)
	if existingUser.ID != 0 {
		return ApiFailed(c, 201, "用户名已存在")
	}

	// Hash password
	hashedPassword, err := util.HashMessage(u.Password)
	if err != nil {
		return ApiFailed(c, 500, "密码加密失败")
	}

	// Create user
	user := model.User{
		Username: u.Username,
		Password: hashedPassword,
		Admin:    u.Admin,
	}

	// Create default config
	cfg := model.Config{
		Zoom:      4,
		MinZom:    2,
		MaxZoom:   10,
		Tolorance: 4,
		IconSize:  10,
		Lon:       110,
		Lat:       32,
	}
	user.Config = cfg

	if err := store.DB.Create(&user).Error; err != nil {
		return ApiFailed(c, 500, err.Error())
	}

	return ApiSuccess(c, user)
}

// @Summary  update user
// @Tags user
// @Security     JWT
// @Param        id    path  int  true  "user id"
// @Accept       json
// @Param        data  body  form.UserAdminUpdateForm true  "data"
// @Success 200 {object} Response
// @Router /api/v1/user/update/{id} [POST]
func UserUpdate(c echo.Context) error {
	id := c.Param("id")

	var user model.User
	if err := store.DB.First(&user, id).Error; err != nil {
		return ApiFailed(c, 404, "用户不存在")
	}

	u := new(form.UserAdminUpdateForm)
	if err := c.Bind(u); err != nil {
		return ApiFailed(c, 400, err.Error())
	}

	// Update password if provided
	if u.Password != "" {
		hashedPassword, err := util.HashMessage(u.Password)
		if err != nil {
			return ApiFailed(c, 500, "密码加密失败")
		}
		user.Password = hashedPassword
	}

	// Update admin role if provided
	if u.Admin != nil {
		user.Admin = *u.Admin
	}

	if err := store.DB.Save(&user).Error; err != nil {
		return ApiFailed(c, 500, err.Error())
	}

	return ApiSuccess(c, user)
}

// @Summary  delete user
// @Tags user
// @Security     JWT
// @Param        id    path  int  true  "user id"
// @Success 200 {object} Response
// @Router /api/v1/user/del/{id} [GET]
func UserDelete(c echo.Context) error {
	id := c.Param("id")

	var user model.User
	if err := store.DB.First(&user, id).Error; err != nil {
		return ApiFailed(c, 404, "用户不存在")
	}

	// Delete user's places and related data
	if err := store.DB.Where("user_id = ?", user.ID).Delete(&model.Photo{}).Error; err != nil {
		util.Logger.Error("删除用户照片失败: " + err.Error())
	}

	if err := store.DB.Where("user_id = ?", user.ID).Delete(&model.Place{}).Error; err != nil {
		util.Logger.Error("删除用户地点失败: " + err.Error())
	}

	if err := store.DB.Where("id = ?", user.ConfigID).Delete(&model.Config{}).Error; err != nil {
		util.Logger.Error("删除用户配置失败: " + err.Error())
	}

	if err := store.DB.Delete(&user).Error; err != nil {
		return ApiFailed(c, 500, err.Error())
	}

	return ApiSuccess(c, nil)
}

type LoginResult struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

// @Summary  user login
// @Tags user
// @Param user body form.UserLoginForm true "UserLoginForm"
// @Success 200 {object} Response
// @Router /api/v1/user/login [POST]
func UserLogin(c echo.Context) error {
	u := new(form.UserLoginForm)
	if err := c.Bind(u); err != nil {
		return err
	}
	config := c.Get("config").(jwter.Config)
	user := model.User{}
	err := store.DB.Where("username = ?", u.Username).Find(&user).Error
	if err == nil && user.Username != "" {
		if util.ValidateHash(user.Password, u.Password) != nil {
			return ApiFailed(c, 201, "password incorrect")
		}
		claims := jwt.MapClaims{
			"id":       user.ID,
			"username": user.Username,
			"admin":    user.Admin,
			"exp":      time.Now().Add(time.Hour * 5).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(config.JWTScrect))
		if err != nil {
			return err
		}
		user.Token = &t
		store.DB.Save(&user)
		return ApiSuccess(c, LoginResult{Token: t, Username: user.Username, Admin: user.Admin})
	} else {
		return ApiFailed(c, 201, "user not exist")
	}

}

// @Summary  user login
// @Tags user
// @Security     JWT
// @Success 200 {object} Response
// @Router /api/v1/user/logout [GET]
func UserLogout(c echo.Context) error {
	user, err := jwter.ParseJWT(c.Get("user"))
	if err == nil && user.Username != "" {
		user.Token = nil
		store.DB.Save(&user)
		return ApiSuccess(c, nil)
	} else {
		return ApiFailed(c, 201, "user not exist")
	}

}

// @Summary  user config
// @Tags user
// @Security     JWT
// @Success 200 {object} Response
// @Router /api/v1/user/config [GET]
func UserConfig(c echo.Context) error {
	user, err := jwter.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	if user.Username != "" {
		return ApiSuccess(c, user.Config)
	} else {
		return ApiFailed(c, 201, "用户名或密码错误")
	}
}

// @Summary update user map config
// @Tags user
// @security JWT
// @Accept       json
// @Produce      json
// @Param        data  body  form.UserConfigUpdateForm true  "data"
// @Success 	 200 {object} Response "Success"
// @Router /api/v1/user/config/update [post]
func UserConfigGetUpate(c echo.Context) error {
	user, err := jwter.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.UserConfigUpdateForm)
	if err := c.Bind(u); err != nil {
		return err
	}
	store.DB.Model(&user).Association("Config").Find(&user.Config)
	if user.Username != "" {
		cfg := user.Config
		if cfg.Title != u.Title {
			cfg.Title = u.Title
		}
		if cfg.AutoCenter != u.AutoCenter {
			cfg.AutoCenter = u.AutoCenter
		}
		if cfg.IconSize != u.IconSize {
			cfg.IconSize = u.IconSize
		}
		if cfg.Zoom != u.Zoom {
			cfg.Zoom = u.Zoom
		}
		if cfg.MaxZoom != u.MaxZoom {
			cfg.MaxZoom = u.MaxZoom
		}
		if cfg.MinZom != u.MinZom {
			cfg.MinZom = u.MinZom
		}
		if cfg.Tolorance != u.Tolorance {
			cfg.Tolorance = u.Tolorance
		}
		if cfg.Lat != u.Lat {
			cfg.Lat = u.Lat
		}
		if cfg.Lon != u.Lon {
			cfg.Lon = u.Lon
		}
		if cfg.Link != u.Link {
			cfg.Link = u.Link
		}
		if cfg.Note != u.Note {
			cfg.Note = u.Note
		}
		store.DB.Save(&cfg)
		return ApiSuccess(c, cfg)
	} else {
		return ApiFailed(c, 201, "无权限")
	}
}

type UserShareResult struct {
	Config   model.Config               `json:"config"`
	Features *geojson.FeatureCollection `json:"features"`
}

// @Summary  user share place
// @Tags 		 user
// @Produce      json
// @Param        username path string    true  "username"
// @Param        group    query string   false "group name"
// @Success 	 200 {object} Response "Success"
// @Router 		 /api/v1/share/{username}  [get]
func UserShare(c echo.Context) error {
	username := c.Param("username")
	groupName := c.QueryParam("group")
	if username == "" {
		return ApiFailed(c, 201, "username is empty")
	}
	user := model.User{}
	err := store.DB.Preload("Config").Preload("Places").Find(&user, "username = ?", username).Error
	if user.ID == 0 || err != nil {
		return ApiFailed(c, 201, "用户不存在")
	}
	fc := geojson.NewFeatureCollection()
	places := user.Places
	for i, v := range places {
		err = store.DB.Preload("Photos").Model(&places[i]).First(&places[i], v.ID).Error
		if err != nil {
			util.Logger.Error(err.Error())
		}
	}
	for _, place := range user.Places {
		if place.Private {
			continue
		}
		// Filter by group if group parameter is provided
		if groupName != "" && place.Group != groupName {
			continue
		}
		pt := orb.Point{place.Lon, place.Lat}
		projectedPT := project.Point(pt, project.WGS84.ToMercator)
		feat := geojson.NewFeature(projectedPT)
		feat.Properties["id"] = place.ID
		feat.Properties["userId"] = place.UserID
		feat.Properties["name"] = place.Name
		feat.Properties["desc"] = place.Desc
		feat.Properties["cover"] = place.Cover
		feat.Properties["lat"] = place.Lat
		feat.Properties["lon"] = place.Lon
		feat.Properties["photos"] = place.Photos
		feat.Properties["date"] = place.Date
		feat.Properties["group"] = place.Group

		fc.Append(feat)
	}
	return ApiSuccess(c, UserShareResult{Config: user.Config, Features: fc})
}
