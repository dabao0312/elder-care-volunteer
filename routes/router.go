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
	}

	return r
}
