package item

import "gorm.io/gorm"

// Item represents an item entity
type Item struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null" binding:"required"`
	Description string `json:"description" gorm:"type:text"`
}
