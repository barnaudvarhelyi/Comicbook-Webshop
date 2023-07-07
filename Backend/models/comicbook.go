package models

type Comicbook struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Price    float32   `json:"price"`
	Category *Category `json:"category"`
}
