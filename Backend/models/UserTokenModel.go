package models

type UserToken struct {
	UserId int64  `json:"userId"`
	Token  string `json:"token"`
}
