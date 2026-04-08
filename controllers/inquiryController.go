package controllers

import (
	"net/http"
	"weavory-backend/config"
	"weavory-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateInquiry(c *gin.Context) {

	var inquiry models.Inquiry

	if err := c.BindJSON(&inquiry); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, err := config.DB.Exec(
		"INSERT INTO inquiries (name,email,contact,message) VALUES ($1,$2,$3,$4)",
		inquiry.Name,
		inquiry.Email,
		inquiry.Contact,
		inquiry.Message,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inquiry sent"})
}