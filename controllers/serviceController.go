package controllers

import (
	"net/http"
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
	"github.com/gin-gonic/gin"

	"path/filepath"
	"time"
	"fmt"
)

func GetServices(c *gin.Context) {

	rows, err := config.DB.Query("SELECT id,title,description,icon FROM services")
	
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	defer rows.Close()

	var services []models.Service

	for rows.Next() {
		var s models.Service

		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.Icon); err != nil {
			utils.Error(c, 500, err.Error())
			return
		}

		services = append(services, s)
	}

	utils.Success(c, services)
}

func CreateService(c *gin.Context) {

	title := c.PostForm("title")
	description := c.PostForm("description")

	file, err := c.FormFile("icon")
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	uploadPath := config.GetEnv("UPLOAD_PATH", "uploads")
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(uploadPath, filename)

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO services (title, description, icon) VALUES ($1,$2,$3)",
		title, description, filePath,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service created",
		"icon":    filePath,
	})
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")

	title := c.PostForm("title")
	description := c.PostForm("description")

	file, _ := c.FormFile("icon")

	var iconPath string

	if file != nil {
		uploadPath := config.GetEnv("UPLOAD_PATH", "uploads")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		iconPath = filepath.Join(uploadPath, filename)
		c.SaveUploadedFile(file, iconPath)
	}

	query := "UPDATE services SET title=$1, description=$2"
	args := []interface{}{title, description}

	if iconPath != "" {
		query += ", icon=$3 WHERE id=$4"
		args = append(args, iconPath, id)
	} else {
		query += " WHERE id=$3"
		args = append(args, id)
	}

	_, err := config.DB.Exec(query, args...)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")

	_, err := config.DB.Exec("DELETE FROM services WHERE id=$1", id)

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}
