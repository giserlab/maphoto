package model

import (
	"time"

	"gorm.io/gorm"
)

type Place struct {
	gorm.Model
	ID      int        `json:"id"`
	UserID  int        `json:"userid"`
	Name    *string    `json:"name"`
	Desc    string     `json:"desc"`
	Cover   string     `json:"cover"`
	Private bool       `json:"private"`
	Group   string     `json:"group" gorm:"default:''"`
	Date    *time.Time `json:"date"`
	Photos  []Photo    `json:"photos"`
	Lon     float64    `json:"lon"`
	Lat     float64    `json:"lat"`
}
