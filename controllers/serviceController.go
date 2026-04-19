package controllers

import (
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
	"github.com/gin-gonic/gin"
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

	iconURL, err := utils.UploadImage(file)
	if err != nil {
		utils.Error(c, 500, "failed upload icon")
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO services (title, description, icon) VALUES ($1,$2,$3)",
		title, description, iconURL,
	)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, "Service created")
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")

	var existing models.Service

	err := config.DB.QueryRow(`
		SELECT title, description, icon FROM services WHERE id=$1
	`, id).Scan(
		&existing.Title,
		&existing.Description,
		&existing.Icon,
	)

	if err != nil {
		utils.Error(c, 404, "Service not found")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" {
		title = existing.Title
	}
	if description == "" {
		description = existing.Description
	}

	icon := existing.Icon

	file, _ := c.FormFile("icon")
	if file != nil {
		url, err := utils.UploadImage(file)
		if err == nil {
			icon = url
		}
	}

	_, err = config.DB.Exec(`
		UPDATE services 
		SET title=$1, description=$2, icon=$3 
		WHERE id=$4
	`, title, description, icon, id)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, "Service updated")
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
