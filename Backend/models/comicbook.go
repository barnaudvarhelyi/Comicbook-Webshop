package models

type Comicbook struct {
	id       int
	name     string
	price    float32
	category *Category
}
