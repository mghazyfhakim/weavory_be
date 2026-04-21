package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"

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

	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	htmlContent := `
		<h2>New Inquiry - Weavory Studio</h2>
		<p><strong>Nama:</strong> ` + inquiry.Name + `</p>
		<p><strong>Email:</strong> ` + inquiry.Email + `</p>
		<p><strong>Contact:</strong> ` + inquiry.Contact + `</p>
		<p><strong>Message:</strong><br/>` + inquiry.Message + `</p>
	`

	params := &resend.SendEmailRequest{
		From:    "Weavory Studio <inquiry@weavorystudio.com>", 
		To:      []string{"weavorystudio@gmail.com"},
		Subject: "New Inquiry from Website",
		Html:    htmlContent,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		println("Email failed:", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inquiry sent successfully",
		"data":    inquiry,
	})
}