package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	PswHash   string `json:"-"`
	CreatedAt string `json:"created_at"`
	Active    bool   `json:"-"`
	VerHash   string `json:"-"`
}
