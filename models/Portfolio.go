package models

type Portfolio struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Material    string `json:"material"`
	TeknikJahit string `json:"teknik_jahit"`
	Finishing   string `json:"finishing"`
	Layanan     string `json:"layanan"`
	ImageURL    string `json:"image_url"` // thumbnail
	Images      []string `json:"images"`  // slider
}
