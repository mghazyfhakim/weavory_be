package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
)

func GetPortfolios(c *gin.Context) {

	limit := c.Query("limit")

	query := "SELECT id,title,description,image_url,material FROM portfolios ORDER BY id DESC"

	if limit != "" {
		var limitVal int
		fmt.Sscanf(limit, "%d", &limitVal)

		if limitVal > 0 {
			query += fmt.Sprintf(" LIMIT %d", limitVal)
		}
	}

	rows, err := config.DB.Query(query)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	defer rows.Close()

	var portfolios []models.Portfolio

	for rows.Next() {
		var s models.Portfolio

		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.ImageURL, &s.Material); err != nil {
			utils.Error(c, 500, err.Error())
			return
		}

		portfolios = append(portfolios, s)
	}

	utils.Success(c, portfolios)
}

func GetPortfolioDetail(c *gin.Context) {
	id := c.Param("id")

	var p models.Portfolio

	err := config.DB.QueryRow(`
		SELECT id, title, description, material, teknik_jahit, finishing, layanan, image_url
		FROM portfolios WHERE id=$1
	`, id).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Material,
		&p.TeknikJahit,
		&p.Finishing,
		&p.Layanan,
		&p.ImageURL,
	)

	if err != nil {
		utils.Error(c, 404, "Portfolio not found")
		return
	}

	rows, err := config.DB.Query(`
		SELECT image_url FROM portfolio_images WHERE portfolio_id=$1
	`, id)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var img string
			rows.Scan(&img)
			p.Images = append(p.Images, img)
		}
	}

	utils.Success(c, p)
}

func CreatePortfolio(c *gin.Context) {
	title := c.PostForm("title")
	material := c.PostForm("material")
	TeknikJahit := c.PostForm("TeknikJahit")
	finishing := c.PostForm("finishing")
	layanan := c.PostForm("layanan")

	if title == "" || material == "" || TeknikJahit == "" || finishing == "" {
		utils.Error(c, 400, "all fields are required")
		return
	}

	description := fmt.Sprintf("%s | %s", TeknikJahit, finishing)

	// 🔥 THUMBNAIL (Cloudinary)
	thumb, err := c.FormFile("thumbnail")
	if err != nil {
		utils.Error(c, 400, "thumbnail is required")
		return
	}

	thumbURL, err := utils.UploadImage(thumb)
	if err != nil {
		utils.Error(c, 500, "failed upload thumbnail")
		return
	}

	var portfolioID int
	err = config.DB.QueryRow(`
		INSERT INTO portfolios 
		(title, description, image_url, material, teknik_jahit, finishing, layanan)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id
	`, title, description, thumbURL, material, TeknikJahit, finishing, layanan).Scan(&portfolioID)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	// 🔥 MULTIPLE IMAGES
	form, _ := c.MultipartForm()
	files := form.File["images"]

	for _, file := range files {

		imageURL, err := utils.UploadImage(file)
		if err != nil {
			continue
		}

		config.DB.Exec(`
			INSERT INTO portfolio_images (portfolio_id, image_url)
			VALUES ($1,$2)
		`, portfolioID, imageURL)
	}

	utils.Success(c, gin.H{
		"id":      portfolioID,
		"message": "Portfolio created with Cloudinary",
	})
}

func UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")

	var existing models.Portfolio

	err := config.DB.QueryRow(`
		SELECT title, material, teknik_jahit, finishing, layanan, image_url
		FROM portfolios WHERE id=$1
	`, id).Scan(
		&existing.Title,
		&existing.Material,
		&existing.TeknikJahit,
		&existing.Finishing,
		&existing.Layanan,
		&existing.ImageURL,
	)

	if err != nil {
		utils.Error(c, 404, "Portfolio not found")
		return
	}

	title := c.PostForm("title")
	material := c.PostForm("material")
	teknikJahit := c.PostForm("TeknikJahit")
	finishing := c.PostForm("finishing")
	layanan := c.PostForm("layanan")

	if title == "" {
		title = existing.Title
	}
	if material == "" {
		material = existing.Material
	}
	if teknikJahit == "" {
		teknikJahit = existing.TeknikJahit
	}
	if finishing == "" {
		finishing = existing.Finishing
	}
	if layanan == "" {
		layanan = existing.Layanan
	}

	description := fmt.Sprintf("%s | %s", teknikJahit, finishing)

	thumb, _ := c.FormFile("thumbnail")
	thumbURL := existing.ImageURL

	if thumb != nil {
		url, err := utils.UploadImage(thumb)
		if err != nil {
			utils.Error(c, 500, "failed upload thumbnail")
			return
		}
		thumbURL = url
	}

	_, err = config.DB.Exec(`
		UPDATE portfolios 
		SET title=$1, description=$2, image_url=$3, material=$4, teknik_jahit=$5, finishing=$6, layanan=$7
		WHERE id=$8
	`, title, description, thumbURL, material, teknikJahit, finishing, layanan, id)

	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["images"]

	if len(files) > 0 {
		config.DB.Exec("DELETE FROM portfolio_images WHERE portfolio_id=$1", id)

		for _, file := range files {
			url, err := utils.UploadImage(file)
			if err != nil {
				continue
			}

			config.DB.Exec(`
				INSERT INTO portfolio_images (portfolio_id, image_url)
				VALUES ($1,$2)
			`, id, url)
		}
	}

	utils.Success(c, gin.H{
		"message": "Portfolio updated (Cloudinary)",
	})
}

func DeletePortfolio(c *gin.Context) {
	id := c.Param("id")

	config.DB.Exec("DELETE FROM portfolios WHERE id=$1", id)

	utils.Success(c, "data deleted")
}
