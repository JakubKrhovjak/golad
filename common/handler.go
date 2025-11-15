package common

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CRUDService defines the interface for CRUD operations
type CRUDService[T any] interface {
	GetAll() ([]*T, error)
	GetByID(id uint) (*T, error)
	Create(entity *T) (*T, error)
	Update(id uint, entity *T) (*T, error)
	Delete(id uint) error
}

// CRUDHandler provides generic CRUD handlers
type CRUDHandler[T any] struct {
	service CRUDService[T]
}

// NewCRUDHandler creates a new generic CRUD handler
func NewCRUDHandler[T any](service CRUDService[T]) *CRUDHandler[T] {
	return &CRUDHandler[T]{
		service: service,
	}
}

// GetAll handles GET requests for all entities
func (h *CRUDHandler[T]) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		entities, err := h.service.GetAll()
		if err != nil {
			RespondError(c, http.StatusInternalServerError, err)
			return
		}
		RespondSuccess(c, http.StatusOK, entities)
	}
}

// GetByID handles GET requests for a single entity by ID
func (h *CRUDHandler[T]) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := ParseID(c)
		if err != nil {
			RespondError(c, http.StatusBadRequest, err)
			return
		}

		entity, err := h.service.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "item not found" {
				RespondError(c, http.StatusNotFound, err)
				return
			}
			RespondError(c, http.StatusInternalServerError, err)
			return
		}
		RespondSuccess(c, http.StatusOK, entity)
	}
}

// Create handles POST requests to create a new entity
func (h *CRUDHandler[T]) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var entity T
		if err := c.ShouldBindJSON(&entity); err != nil {
			RespondError(c, http.StatusBadRequest, err)
			return
		}

		created, err := h.service.Create(&entity)
		if err != nil {
			RespondError(c, http.StatusInternalServerError, err)
			return
		}
		RespondSuccess(c, http.StatusCreated, created)
	}
}

// Update handles PUT requests to update an entity
func (h *CRUDHandler[T]) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := ParseID(c)
		if err != nil {
			RespondError(c, http.StatusBadRequest, err)
			return
		}

		var entity T
		if err := c.ShouldBindJSON(&entity); err != nil {
			RespondError(c, http.StatusBadRequest, err)
			return
		}

		updated, err := h.service.Update(id, &entity)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "item not found" {
				RespondError(c, http.StatusNotFound, err)
				return
			}
			RespondError(c, http.StatusInternalServerError, err)
			return
		}
		RespondSuccess(c, http.StatusOK, updated)
	}
}

// Delete handles DELETE requests to delete an entity
func (h *CRUDHandler[T]) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := ParseID(c)
		if err != nil {
			RespondError(c, http.StatusBadRequest, err)
			return
		}

		if err := h.service.Delete(id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "item not found" {
				RespondError(c, http.StatusNotFound, err)
				return
			}
			RespondError(c, http.StatusInternalServerError, err)
			return
		}
		RespondSuccess(c, http.StatusOK, gin.H{"message": "Entity deleted successfully"})
	}
}

// ParseID extracts and parses ID from URL parameter
func ParseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID")
	}
	return uint(id), nil
}

// RespondSuccess sends a successful JSON response
func RespondSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// RespondError sends an error JSON response
func RespondError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success message response
type SuccessResponse struct {
	Message string `json:"message"`
}
