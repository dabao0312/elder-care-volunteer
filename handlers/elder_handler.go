package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"elder-care-volunteer/models"
)

type ElderRegisterRequest struct {
	Phone         string  `json:"phone"`
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	GuardianPhone string  `json:"guardian_phone"`
}

func RegisterElder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ElderRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {

			user := models.User{
				Phone:  req.Phone,
				Name:   req.Name,
				Role:   "elder",
				Status: "active",
			}

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			elder := models.Elder{
				UserID:        user.ID,
				Address:       req.Address,
				Latitude:      req.Latitude,
				Longitude:     req.Longitude,
				GuardianPhone: req.GuardianPhone,
			}

			if err := tx.Create(&elder).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "elder registered successfully"})
	}
}
