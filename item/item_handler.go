package item

import (
	"awesomeProject2/common"
)

// Handler handles HTTP requests for items
// Embeds generic CRUDHandler to get all CRUD operations automatically
type Handler struct {
	*common.CRUDHandler[Item]
	service *ItemService
}

// NewHandler creates a new Handler instance
func NewHandler(service *ItemService) *Handler {
	return &Handler{
		CRUDHandler: common.NewCRUDHandler[Item](service),
		service:     service,
	}
}

// You can add custom handler methods here if needed
// Example:
// func (h *Handler) SearchByName() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         name := c.Query("name")
//         items, err := h.service.SearchByName(name)
//         if err != nil {
//             common.RespondError(c, http.StatusInternalServerError, err)
//             return
//         }
//         common.RespondSuccess(c, http.StatusOK, items)
//     }
// }
