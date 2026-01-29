package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"elder-care-volunteer/models"
)

type VolunteerRegisterRequest struct {
	Phone     string  `json:"phone" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func RegisterVolunteer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VolunteerRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {

			user := models.User{
				Phone:  req.Phone,
				Name:   req.Name,
				Role:   "volunteer",
				Status: "active",
			}

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			volunteer := models.VolunteerProfile{
				UserID:    user.ID,
				Address:   req.Address,
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
			}

			if err := tx.Create(&volunteer).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "volunteer registered successfully"})
	}
}
