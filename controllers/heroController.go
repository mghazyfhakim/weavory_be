package controllers

import (
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
	"github.com/gin-gonic/gin"
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
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, hero)
}

func UpdateHero(c *gin.Context) {

	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	description := c.PostForm("description")

	var existing models.Hero

	err := config.DB.QueryRow(`
		SELECT title, subtitle, description, image_url
		FROM hero WHERE id=1
	`).Scan(
		&existing.Title,
		&existing.Subtitle,
		&existing.Description,
		&existing.ImageURL,
	)

	if err != nil {
		utils.Error(c, 404, "Hero not found")
		return
	}

	if title == "" {
		title = existing.Title
	}
	if subtitle == "" {
		subtitle = existing.Subtitle
	}
	if description == "" {
		description = existing.Description
	}

	// 🔥 CLOUDINARY
	file, _ := c.FormFile("image_url")
	imageURL := existing.ImageURL

	if file != nil {
		url, err := utils.UploadImage(file)
		if err != nil {
			utils.Error(c, 500, "failed upload image")
			return
		}
		imageURL = url
	}

	_, err = config.DB.Exec(`
		UPDATE hero 
		SET title=$1, subtitle=$2, description=$3, image_url=$4
		WHERE id=1
	`, title, subtitle, description, imageURL)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, "Hero updated")
}