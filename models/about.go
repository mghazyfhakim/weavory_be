package models

type About struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	ImageURL    string   `json:"image_url"`

	Profile     string   `json:"profile"`
	Vision      string   `json:"vision"`
	Mission     []string `json:"mission"`
}	