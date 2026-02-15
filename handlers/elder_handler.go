package handlers

import (
	"log"
	"math"
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

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// 将角度转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// 使用 Haversine 公式计算距离 (单位: 公里)
	a := math.Sin((lat2Rad-lat1Rad)/2)*math.Sin((lat2Rad-lat1Rad)/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin((lon2Rad-lon1Rad)/2)*math.Sin((lon2Rad-lon1Rad)/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 地球半径约为 6371 公里
	return 6371 * c
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
func ElderNoReply(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		elderID := c.Param("id")

		// 1️⃣ 查老人档案是否存在
		var elder models.Elder
		if err := db.First(&elder, elderID).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "elder not found",
			})
			return
		}

		// 获取关联的用户信息以取得姓名
		var user models.User
		if err := db.Where("id = ?", elder.UserID).First(&user).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "elder's user info not found",
			})
			return
		}

		// // 2️⃣ 找一个可用志愿者（先随便找一个）
		// var volunteer models.VolunteerProfile
		// if err := db.Where("available = ?", 1).First(&volunteer).Error; err != nil {
		// 	c.JSON(404, gin.H{
		// 		"error": "no available volunteer",
		// 	})
		// 	return
		// }

		var volunteers []models.VolunteerProfile
		if err := db.Where("available = ?", 1).Find(&volunteers).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(volunteers) == 0 {
			c.JSON(404, gin.H{"error": "no available volunteer"})
			return
		}

		// 找最近的志愿者
		minDistance := math.MaxFloat64
		var nearest models.VolunteerProfile

		for _, v := range volunteers {
			d := distance(
				elder.Latitude, elder.Longitude,
				v.Latitude, v.Longitude,
			)

			if d < minDistance {
				minDistance = d
				nearest = v
			}
		}

		// 3️⃣ 模拟"通知"
		var volunteerUser models.User
		if err := db.Where("id = ?", nearest.UserID).First(&volunteerUser).Error; err != nil {
			c.JSON(404, gin.H{
				"error": "volunteer's user info not found",
			})
			return
		}

		// log.Printf(
		// 	"[ALERT] elder %d (%s) no reply, notify volunteer %d (%s)",
		// 	elder.ID, user.Name,
		// 	nearest.ID, volunteerUser.Name,
		// )

		log.Printf(
			"[ALERT] elder %d (%s) no reply, notify volunteer %d (%s), distance=%.2fkm",
			elder.ID, user.Name,
			nearest.ID, volunteerUser.Name,
			minDistance,
		)

		c.JSON(http.StatusOK, gin.H{
			"message":      "nearest volunteer notified (simulated)",
			"elder_id":     elder.ID,
			"volunteer_id": nearest.ID,
			"distance_km":  minDistance,
		})

	}
}
