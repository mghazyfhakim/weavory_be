package controllers

import (
	"weavory-backend/config"
	"weavory-backend/models"

	"github.com/gin-gonic/gin"
	"path/filepath"
)

func GetHero(c *gin.Context) {

	var hero models.Hero

	err := config.DB.QueryRow(
		"SELECT id,title,subtitle,description,image_url FROM hero LIMIT 1",
	).Scan(
		&hero.ID,
		&hero.Title,
		&hero.Subtitle,
		&hero.Description,
		&hero.ImageURL,
	)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, hero)
}



func UpdateHero(c *gin.Context) {

	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	description := c.PostForm("description")

	file, _ := c.FormFile("image_url")

	var imagePath string

	if file != nil {
		filename := filepath.Base(file.Filename)
		imagePath = filepath.Join("uploads", filename)
		c.SaveUploadedFile(file, imagePath)
	}

	query := "UPDATE hero SET title=$1, subtitle=$2, description=$3, image_url=$4 WHERE id=1"

	args := []interface{}{title, subtitle, description, imagePath}

	_, err := config.DB.Exec(query, args...)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "Hero updated"})
}