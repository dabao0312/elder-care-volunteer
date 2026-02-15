package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"elder-care-volunteer/models"
)

func ElderReply(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		elderID := c.Param("id")

		// 1️⃣ 更新未处理告警为 resolved
		result := db.Model(&models.AlertRecord{}).
			Where("elder_id = ? AND status = ?", elderID, "pending").
			Update("status", "resolved")

		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message":  "elder replied, alerts resolved",
			"elder_id": elderID,
		})
	}
}
