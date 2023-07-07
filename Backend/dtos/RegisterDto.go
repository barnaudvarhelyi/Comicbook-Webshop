package dtos

type RegisterDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ConfPassword string `json:"conf_password"`
	Email        string `json:"email"`
}
