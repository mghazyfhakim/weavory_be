package models

type About struct {
	ID          int    `json:"id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	ImageURL    string `json:"image_url" form:"image_url"`
}