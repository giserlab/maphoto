package form

type PlaceUpdateForm struct {
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
	Group string  `json:"group"`
	Lon   float64 `json:"lon"`
	Lat   float64 `json:"lat"`
	Cover string  `json:"cover"`
}
