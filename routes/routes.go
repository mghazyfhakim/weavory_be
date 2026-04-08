package routes

import (
	"weavory-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		// services
		api.GET("/services", controllers.GetServices)
		api.POST("/services", controllers.CreateService)
		api.PUT("/services/:id", controllers.UpdateService)
		api.DELETE("/services/:id", controllers.DeleteService)

		// portfolio
		api.GET("/portfolios", controllers.GetPortfolios)
		api.POST("/portfolios", controllers.CreatePortfolio)
		api.PUT("/portfolios/:id", controllers.UpdatePortfolio)
		api.DELETE("/portfolios/:id", controllers.DeletePortfolio)

		// inquiry
		api.POST("/inquiry", controllers.CreateInquiry)

		api.GET("/hero", controllers.GetHero)
		api.PUT("/hero", controllers.UpdateHero)

		api.GET("/about", controllers.GetAbout)
		api.PUT("/about", controllers.UpdateAbout)
	}
}
