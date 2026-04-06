package model

type Config struct {
	Id         int     `json:"id"`
	UserId     int     `json:"userId"`
	Title      string  `json:"title"`
	Link       string  `json:"link"`
	IconSize   int     `json:"iconSize"`
	Lon        float32 `json:"lon"`
	Lat        float32 `json:"lat"`
	Zoom       int     `json:"zoom"`
	MaxZoom    int     `json:"maxZoom"`
	MinZom     int     `json:"minZoom"`
	Tolorance  float32 `json:"tolorance"`
	AutoCenter bool    `json:"autoCenter"`
	Note       string  `json:"note"`
}
