package controllers

import (
	"net/http"
	"weavory-backend/config"
	"weavory-backend/models"

	"github.com/gin-gonic/gin"
	"path/filepath"
)

func GetPortfolios(c *gin.Context) {

	rows, err := config.DB.Query("SELECT id,title,description,image_url FROM portfolios")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var portfolios []models.Portfolio

	for rows.Next() {
		var s models.Portfolio

		rows.Scan(&s.ID, &s.Title, &s.Description, &s.ImageURL)

		portfolios = append(portfolios, s)
	}

	c.JSON(http.StatusOK, portfolios)
}

func CreatePortfolio(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")

	file, err := c.FormFile("image_url")
	if err != nil {
		c.JSON(400, gin.H{"error": "image_url is required"})
		return
	}

	filename := filepath.Base(file.Filename)
	filePath := filepath.Join("uploads", filename)

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = config.DB.Exec(
		"INSERT INTO portfolios (title, description, image_url) VALUES ($1,$2,$3)",
		title, description, filePath,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portfolios created",
		"image_url":    filePath,
	})
}

func UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")

	title := c.PostForm("title")
	description := c.PostForm("description")

	file, _ := c.FormFile("image_url")

	var imagePath string

	if file != nil {
		filename := filepath.Base(file.Filename)
		imagePath = filepath.Join("uploads", filename)
		c.SaveUploadedFile(file, imagePath)
	}

	query := "UPDATE portfolios SET title=$1, description=$2"
	args := []interface{}{title, description}

	if imagePath != "" {
		query += ", image_url=$3 WHERE id=$4"
		args = append(args, imagePath, id)
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

func DeletePortfolio(c *gin.Context) {
	id := c.Param("id")

	config.DB.Exec("DELETE FROM portfolios WHERE id=$1", id)

	c.JSON(200, gin.H{"message": "Deleted"})
}