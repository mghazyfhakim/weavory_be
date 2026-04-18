package models

type Inquiry struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Contact string `json:"contact" binding:"required"`
	Message string `json:"message" binding:"required"`
}