package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"elder-care-volunteer/handlers"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/elders/register", handlers.RegisterElder(db))
		api.POST("/volunteers/register", handlers.RegisterVolunteer(db))
		api.GET("/volunteers/available", handlers.ListAvailableVolunteers(db))
		api.POST("/elders/:id/no-reply", handlers.ElderNoReply(db))
		api.POST("/elders/:id/reply", handlers.ElderReply(db))

	}

	return r
}
