package models

import "time"

type VolunteerProfile struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	UserID    int64     `gorm:"column:user_id"`
	Address   string    `gorm:"column:address"`
	Latitude  float64   `gorm:"column:latitude"`
	Longitude float64   `gorm:"column:longitude"`
	Available int       `gorm:"column:available"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (VolunteerProfile) TableName() string {
	return "volunteer_profile"
}
