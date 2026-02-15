package models

import "time"

type AlertRecord struct {
	ID          int64 `gorm:"primaryKey"`
	ElderID     int64
	VolunteerID int64
	DistanceKM  float64
	Status      string
	CreatedAt   time.Time
	NotifiedAt  *time.Time
	AcceptedAt  *time.Time
	HandledAt   *time.Time
	ClosedAt    *time.Time
}

func (AlertRecord) TableName() string {
	return "alert_record"
}
