package controllers

import (
	"github.com/gin-gonic/gin"
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
)

func GetMaterials(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, title, description FROM materials
	`)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	defer rows.Close()

	var materials []models.Material

	for rows.Next() {
		var m models.Material
		rows.Scan(&m.ID, &m.Title, &m.Description)
		materials = append(materials, m)
	}

	utils.Success(c, materials)
}