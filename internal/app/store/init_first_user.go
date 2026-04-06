package store

import (
	"maphoto/internal/app/model"
	"maphoto/internal/util"
)

func InitAdmin(username string, password string) {
	user := model.User{}
	DB.Where("username = ?", username).Find(&user)
	toke, _ := util.GenerateRandomKey(12)
	pass, err := util.HashMessage(password)
	if err != nil {
		panic("init admin error：" + err.Error())
	}
	if user.Username == "" || user.Password == "" {
		user.Token = &toke
		user.Username = username
		user.Admin = true
		user.Password = pass
		cfg := model.Config{
			UserId:    user.ID,
			Zoom:      4,
			MinZom:    2,
			MaxZoom:   10,
			Tolorance: 4,
			IconSize:  10,
			Lon:       110,
			Lat:       32,
		}
		user.Config = cfg
		DB.Create(&user)
	}
}
