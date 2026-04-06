package api

import (
	"fmt"
	"maphoto/internal/app/form"
	"maphoto/internal/app/model"
	"maphoto/internal/app/store"
	"maphoto/internal/app/web/jwt"
	"maphoto/internal/util"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/project"
)

func FindPlace(user *model.User, id string) *model.Place {
	tID, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	place := &model.Place{}
	// 直接查询数据库，确保地点属于当前用户
	err = store.DB.Preload("Photos").Where("id = ? AND user_id = ?", tID, user.ID).First(place).Error
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	return place
}

// @Summary get place
// @Tags place
// @Accept json
// @security JWT
// @Success 200 {object} Response "Success"
// @Router /api/v1/place/all [get]
func PlaceAll(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	places := user.Places

	for i, v := range places {
		err = store.DB.Preload("Photos").Model(&places[i]).First(&places[i], v.ID).Error
		if err != nil {
			util.Logger.Error(err.Error())
		}
	}
	// 按创建时间倒序排列（新的在前）
	sort.Slice(places, func(i, j int) bool {
		return places[i].ID > places[j].ID
	})
	return ApiSuccess(c, places)
}

// @Summary add place
// @Tags place
// @security JWT
// @Accept       json
// @Produce      json
// @Param        data	body	form.PlaceAddForm	true  "PlaceAddForm"
// @Success 	 200 {object} Response "Success"
// @Router /api/v1/place/add [post]
func PlaceAdd(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.PlaceAddForm)
	if err := c.Bind(u); err != nil {
		return err
	}

	if user.Username == "" {
		return ApiFailed(c, 201, "用户不存在")

	} else {

		time := time.Now()
		pt := model.Place{
			UserID: user.ID,
			Name:   &u.Name,
			Desc:   u.Desc,
			Cover:  u.Cover,
			Group:  u.Group,
			Date:   &time,
			Lat:    u.Lat,
			Lon:    u.Lon,
		}
		store.DB.Save(&pt)
		if len(u.Photos) > 0 {
			for _, v := range u.Photos {
				p := model.Photo{
					PlaceID: pt.ID,
					Url:     v,
				}
				store.DB.Model(&pt).Association("Photos").Append(&p)
			}
		}
		store.DB.Model(&user).Association("Places").Append(&pt)
		store.DB.Save(&user)
		return ApiSuccess(c, pt)
	}

}

// @Summary delete place
// @Tags place
// @security JWT
// @Accept       json
// @Produce      json
// @Param        id		path	int    true  "id"
// @Success 	 200	{object}	Response	"Success"
// @Router /api/v1/place/del/{id} [get]
func PlaceDel(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	id := c.Param("id")
	pt := FindPlace(user, id)
	if pt == nil {
		return ApiFailed(c, 201, "地点不存在")
	} else {
		// 先删除地点的照片
		if err := store.DB.Where("place_id = ?", pt.ID).Delete(&model.Photo{}).Error; err != nil {
			util.Logger.Error("删除地点照片失败: " + err.Error())
			return ApiFailed(c, 201, "删除照片失败")
		}
		// 删除地点记录
		if err := store.DB.Delete(pt).Error; err != nil {
			util.Logger.Error("删除地点失败: " + err.Error())
			return ApiFailed(c, 201, "删除地点失败")
		}
		return ApiSuccess(c, pt)
	}
}

// @Summary 	update props of  place
// @Tags 		 place
// @security 	 JWT
// @Accept       json
// @Produce      json
// @Param        id		path		int    true  "id"
// @Param        data	body		form.PlaceUpdateForm	true  "data"
// @Success 	 200	{object}	Response "Success"
// @Router 		 /api/v1/place/update/{id} [post]
func PlaceUpdate(c echo.Context) error {
	id := c.Param("id")
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.PlaceUpdateForm)
	if err := c.Bind(u); err != nil {
		return err
	}

	pt := model.Place{}
	err = store.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&pt).Error

	if err != nil {
		return ApiFailed(c, 201, "地点不存在")
	} else {
		if u.Name != "" {
			pt.Name = &u.Name
		}
		if u.Desc != "" {
			pt.Desc = u.Desc
		}
		pt.Group = u.Group
		if u.Lon != 0 {
			pt.Lon = u.Lon
		}
		if u.Lat != 0 {
			pt.Lat = u.Lat
		}
		pt.Cover = u.Cover
		store.DB.Save(&pt)
		// 重新查询以获取完整数据
		store.DB.Preload("Photos").First(&pt, pt.ID)
		return ApiSuccess(c, pt)
	}
}

// @Summary place change cover
// @Tags place
// @security JWT
// @Accept       json
// @Produce      json
// @Param        data   body form.PlaceCoverForm true "PlaceCoverForm"
// @Success 	 200 {object} Response "Success"
// @Router /api/v1/place/cover [post]
func PlaceCover(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.PlaceCoverForm)
	if err := c.Bind(u); err != nil {
		return err
	}
	pt := FindPlace(user, fmt.Sprintf("%v", u.Id))
	if pt == nil {
		return ApiFailed(c, 1, "地点不存在")
	} else {
		pt.Cover = u.Url
		store.DB.Save(&pt)
		return ApiSuccess(c, pt)
	}
}

// @Summary place add picture
// @Tags place
// @security JWT
// @Accept       json
// @Produce      json
// @Param	data	body	form.PlaceAttachForm true "PlaceAttachForm"
// @Success 	 200	{object}	Response	"Success"
// @Router /api/v1/place/pic/add [post]
func PlacePicAdd(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.PlaceAttachForm)
	if err := c.Bind(u); err != nil {
		return err
	}
	pt := FindPlace(user, fmt.Sprintf("%v", u.Id))

	if pt == nil {
		return ApiFailed(c, 201, "地点不存在")
	} else {
		photoExist := false
		for _, p := range pt.Photos {
			if p.Url == u.Url {
				photoExist = true
				break
			}
		}
		if photoExist {
			return ApiFailed(c, 201, "图片已存在")
		}

		photo := model.Photo{
			PlaceID: pt.ID,
			Url:     u.Url,
		}
		store.DB.Model(&pt).Association("Photos").Append(&photo)
		// 重新查询以获取最新的 photos
		store.DB.Preload("Photos").First(&pt, pt.ID)
		return ApiSuccess(c, pt)
	}
}

// @Summary  delete place pic
// @Tags 		 place
// @security 	 JWT
// @Accept       json
// @Produce      json
// @Param        data	body	form.PlacePicDelPost    true  "data"
// @Success 	 200	{object}	Response	"Success"
// @Router 		 /api/v1/place/pic/del  [post]
func PlacePicDel(c echo.Context) error {
	user, err := jwt.ParseJWT(c.Get("user"))
	if err != nil {
		return ApiFailed(c, 201, err.Error())
	}
	u := new(form.PlacePicDelPost)
	if err := c.Bind(u); err != nil {
		return err
	}

	pt := FindPlace(user, fmt.Sprintf("%v", u.Id))
	if pt == nil {
		return ApiFailed(c, 201, "地点不存在")
	} else {
		photos := pt.Photos
		for i, v := range photos {
			if v.Url == u.Url {
				store.DB.Model(&pt).Association("Photos").Unscoped().Delete(&photos[i])
				break
			}

		}
		store.DB.Save(&pt)
		// 重新查询以获取最新的 photos
		store.DB.Preload("Photos").First(&pt, pt.ID)
		return ApiSuccess(c, pt)
	}
}

// @Summary   place init
// @Tags 		 place
// @Produce      json
// @security 	 JWT
// @Param        username path string    true  "username"
// @Success 	 200 {object} Response "Success"
// @Router 		 /api/v1/place/init/{username} [get]
func PlaceInit(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return ApiFailed(c, 201, "username is empty")
	}
	user := model.User{}
	err := store.DB.Preload("Config").Preload("Places").Find(&user, "username = ?", username).Error
	if user.ID == 0 || err != nil {
		return ApiFailed(c, 201, "用户不存在")
	}
	fc := geojson.NewFeatureCollection()
	for _, place := range user.Places {
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
		feat.Properties["date"] = place.Date

		fc.Append(feat)
	}
	return ApiSuccess(c, UserShareResult{Config: user.Config, Features: fc})
}
