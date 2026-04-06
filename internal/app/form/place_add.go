package form

type PlaceAddForm struct {
	Name   string   `json:"name"`
	Desc   string   `json:"desc"`
	Lon    float64  `json:"lon"`
	Lat    float64  `json:"lat"`
	Cover  string   `json:"cover"`
	Group  string   `json:"group"`
	Photos []string `json:"photos"`
}
