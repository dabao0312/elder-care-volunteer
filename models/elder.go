package models

import "time"

type Elder struct {
	ID            int64     `gorm:"primaryKey;column:id"`
	UserID        int64     `gorm:"column:user_id"`
	Address       string    `gorm:"column:address"`
	Latitude      float64   `gorm:"column:latitude"`
	Longitude     float64   `gorm:"column:longitude"`
	GuardianPhone string    `gorm:"column:guardian_phone"`
	LastReplyTime *time.Time `gorm:"column:last_reply_time"`
	CreatedAt     time.Time `gorm:"column:create_at"`
}

//
func (Elder) TableName() string {
	return "elder_profile"
}
