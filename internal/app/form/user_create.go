package form

type UserCreateForm struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
	Admin    bool   `json:"admin"`
}
