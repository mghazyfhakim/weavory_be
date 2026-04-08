package controllers

import (
	"weavory-backend/config"
	"weavory-backend/models"

	"github.com/gin-gonic/gin"
	"path/filepath"
)

func GetAbout(c *gin.Context) {

	var about models.About

	err := config.DB.QueryRow(
		"SELECT id,title,description,image_url FROM about LIMIT 1",
	).Scan(
		&about.ID,
		&about.Title,
		&about.Description,
		&about.ImageURL,
	)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, about)
}

func UpdateAbout(c *gin.Context) {

	title := c.PostForm("title")
	description := c.PostForm("description")

	file, _ := c.FormFile("image_url")

	var imagePath string

	if file != nil {
		filename := filepath.Base(file.Filename)
		imagePath = filepath.Join("uploads", filename)
		c.SaveUploadedFile(file, imagePath)
	}

	query := "UPDATE about SET title=$1, description=$2, image_url=$3 WHERE id=1"
	args := []interface{}{title, description, imagePath}

	_, err := config.DB.Exec(query, args...)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "About updated"})
}