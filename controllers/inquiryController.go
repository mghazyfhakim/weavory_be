package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
)

func CreateInquiry(c *gin.Context) {

	var inquiry models.Inquiry

	if err := c.BindJSON(&inquiry); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if inquiry.Name == "" || inquiry.Email == "" || inquiry.Contact == "" {
		utils.Error(c, 400, "Semua field wajib diisi")
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
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inquiry sent successfully",
		"data":    inquiry,
	})
}
