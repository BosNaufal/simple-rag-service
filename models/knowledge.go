package models

import (
	"database/sql"
	"time"
)

// equals
type Knowledge struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	Embedding sql.NullString `gorm:"->:false;<-:create"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp"`
	Distance  float64        `gorm:"->:create"`
}

// TableName overrides the table name used by User to `profiles`
func (Knowledge) TableName() string {
	return "note_chunks"
}
