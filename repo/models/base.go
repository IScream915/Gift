package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
}
