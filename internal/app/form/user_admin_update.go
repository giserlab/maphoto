package form

type UserAdminUpdateForm struct {
	Password string `json:"password,omitempty"`
	Admin    *bool  `json:"admin,omitempty"`
}
