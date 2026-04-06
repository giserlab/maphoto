package model

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	ID      int    `json:"id"`
	PlaceID int    `json:"placeId"`
	Url     string `json:"url"`
}
