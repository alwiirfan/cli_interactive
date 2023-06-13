package entity

import "time"

type Wikis struct {
	ID          int64     `gorm:"primaryKey;column:wikis_id" json:"wikis_id"`
	Topic       string    `gorm:"type:varchar(99)" json:"topic"`
	Description string    `gorm:"type:text" json:"description"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}
