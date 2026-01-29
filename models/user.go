package models

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Phone     string    `gorm:"column:phone"`
	Name      string    `gorm:"column:name"`
	Role      string    `gorm:"column:role"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:create_at"`
}

// 明确指定表名
func (User) TableName() string {
	return "user"
}
