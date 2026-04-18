package controllers

import (
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"time"
	"fmt"
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

	file, _ := c.FormFile("image_url")

	var imagePath string

	if file != nil {
		uploadPath := config.GetEnv("UPLOAD_PATH", "uploads")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		imagePath = filepath.Join(uploadPath, filename)
		c.SaveUploadedFile(file, imagePath)
	}

	query := "UPDATE about SET title=$1, subtitle=$2, description=$3"
	args := []interface{}{title, subtitle, description}

	if imagePath != "" {
	query += ", image_url=4 WHERE id=1"
	args = append(args, imagePath)
	} else {
	query += " WHERE id=1"
	}

	_, err := config.DB.Exec(query, args...)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, "Hero updated")
}