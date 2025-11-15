package item

import (
	"errors"

	"gorm.io/gorm"
)

// ItemService handles business logic for items
type ItemService struct {
	db *gorm.DB
}

// NewItemService creates a new ItemService instance
func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{
		db: db,
	}
}

// GetAll returns all items
func (s *ItemService) GetAll() ([]*Item, error) {
	var items []*Item
	result := s.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

// GetByID returns an item by ID
func (s *ItemService) GetByID(id uint) (*Item, error) {
	var item Item
	result := s.db.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}
	return &item, nil
}

// Create creates a new item
func (s *ItemService) Create(item *Item) (*Item, error) {
	result := s.db.Create(item)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

// Update updates an existing item
func (s *ItemService) Update(id uint, updatedItem *Item) (*Item, error) {
	var item Item
	result := s.db.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, result.Error
	}

	item.Name = updatedItem.Name
	item.Description = updatedItem.Description

	result = s.db.Save(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

// Delete deletes an item by ID
func (s *ItemService) Delete(id uint) error {
	result := s.db.Delete(&Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}
