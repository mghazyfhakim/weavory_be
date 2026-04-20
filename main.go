package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"weavory-backend/config"
	"weavory-backend/routes"
	"weavory-backend/utils"
)

func main() {
	utils.InitCloudinary()

	config.ConnectDB()

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://weavorystudio.com",
			"https://www.weavorystudio.com",
			"http://localhost:3000",

		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
