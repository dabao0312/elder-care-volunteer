package services

import (
	"log"
	"math"
	"time"

	"gorm.io/gorm"

	"elder-care-volunteer/models"
)

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	a := math.Sin((lat2Rad-lat1Rad)/2)*math.Sin((lat2Rad-lat1Rad)/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin((lon2Rad-lon1Rad)/2)*math.Sin((lon2Rad-lon1Rad)/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return 6371 * c
}

func AlertNearestVolunteer(db *gorm.DB, elder models.Elder) {

	//  🛑 防重复告警
	// var existing models.AlertRecord
	// err := db.Where("elder_id = ? AND status = ?", elder.ID, "pending").
	// 	First(&existing).Error

	// if err == nil {
	// 	log.Printf("[SKIP] elder %d already has pending alert", elder.ID)
	// 	return
	// }
	var count int64

	db.Model(&models.AlertRecord{}).
		Where("elder_id = ? AND status = ?", elder.ID, "pending").
		Count(&count)

	if count > 0 {
		log.Printf("[SKIP] elder %d already has pending alert", elder.ID)
		return
	}

	var volunteers []models.VolunteerProfile
	db.Where("available = ?", 1).Find(&volunteers)

	if len(volunteers) == 0 {
		log.Println("[TASK] no available volunteer")
		return
	}

	minDistance := math.MaxFloat64
	var nearest models.VolunteerProfile

	for _, v := range volunteers {
		d := distance(elder.Latitude, elder.Longitude, v.Latitude, v.Longitude)
		if d < minDistance {
			minDistance = d
			nearest = v
		}
	}

	log.Printf(
		"[AUTO ALERT] elder %d notify volunteer %d distance %.2fkm",
		elder.ID,
		nearest.ID,
		minDistance,
	)

	// record := models.AlertRecord{
	// 	ElderID:     elder.ID,
	// 	VolunteerID: nearest.ID,
	// 	DistanceKM:  minDistance,
	// 	Status:      "pending",
	// }
	now := time.Now()

	record := models.AlertRecord{
		ElderID:     elder.ID,
		VolunteerID: nearest.ID,
		DistanceKM:  minDistance,
		Status:      "pending",
		CreatedAt:   now,
		NotifiedAt:  &now,
	}

	if err := db.Create(&record).Error; err != nil {
		log.Println("[ALERT] save record failed:", err)
	}

}
