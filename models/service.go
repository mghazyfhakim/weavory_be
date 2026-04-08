package models

type Service struct {
	ID          int    `json:"id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Icon        string `json:"icon" form:"icon"`
}
