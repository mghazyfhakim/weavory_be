package models

type Hero struct {
	ID          int    `json:"id"`
	Title       string `json:"title" form:"title"`
	Subtitle    string `json:"subtitle" form:"subtitle"`
	Description string `json:"description" form:"description"`
	ImageURL    string `json:"image_url" form:"image_url"`
}