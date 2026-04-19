package controllers

import (
	"encoding/json"
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
	"github.com/gin-gonic/gin"
	"fmt"
	
)

func GetAbout(c *gin.Context) {

	var about models.About
	var missionStr string

	err := config.DB.QueryRow(`
		SELECT id,title,description,image_url,profile,vision,mission 
		FROM about LIMIT 1
	`).Scan(
		&about.ID,
		&about.Title,
		&about.Description,
		&about.ImageURL,
		&about.Profile,
		&about.Vision,
		&missionStr,
	)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	if missionStr != "" {
		json.Unmarshal([]byte(missionStr), &about.Mission)
	}

	utils.Success(c, about)
}

func UpdateAbout(c *gin.Context) {

	var (
		query = "UPDATE about SET "
		args  []interface{}
		index = 1
	)

	title := c.PostForm("title")
	description := c.PostForm("description")
	profile := c.PostForm("profile")
	vision := c.PostForm("vision")
	missionStr := c.PostForm("mission")

	if title != "" {
		query += fmt.Sprintf("title=$%d,", index)
		args = append(args, title)
		index++
	}

	if description != "" {
		query += fmt.Sprintf("description=$%d,", index)
		args = append(args, description)
		index++
	}

	if profile != "" {
		query += fmt.Sprintf("profile=$%d,", index)
		args = append(args, profile)
		index++
	}

	if vision != "" {
		query += fmt.Sprintf("vision=$%d,", index)
		args = append(args, vision)
		index++
	}

	if missionStr != "" {
		var missions []string

		err := json.Unmarshal([]byte(missionStr), &missions)
		if err != nil {
			utils.Error(c, 400, "invalid mission format")
			return
		}

		missionJSON, _ := json.Marshal(missions)

		query += fmt.Sprintf("mission=$%d,", index)
		args = append(args, string(missionJSON))
		index++
	}

	file, _ := c.FormFile("image_url")
	if file != nil {
		url, err := utils.UploadImage(file)
		if err == nil {
			query += fmt.Sprintf("image_url=$%d,", index)
			args = append(args, url)
			index++
		}
	}

	if len(args) == 0 {
		utils.Error(c, 400, "no data to update")
		return
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id=$%d", index)
	args = append(args, 1)

	_, err := config.DB.Exec(query, args...)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, "About updated")
}