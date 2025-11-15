# Generic CRUD Handler

Tento balíček obsahuje generické utility pro vytváření CRUD API s minimálním boilerplate kódem.

## Použití

### 1. Definujte svůj model a service

```go
// Model
type Item struct {
    gorm.Model
    Name        string `json:"name"`
    Description string `json:"description"`
}

// Service musí implementovat CRUDService interface
type ItemService struct {
    db *gorm.DB
}

func (s *ItemService) GetAll() ([]*Item, error) { ... }
func (s *ItemService) GetByID(id uint) (*Item, error) { ... }
func (s *ItemService) Create(entity *Item) (*Item, error) { ... }
func (s *ItemService) Update(id uint, entity *Item) (*Item, error) { ... }
func (s *ItemService) Delete(id uint) error { ... }
```

### 2. Vytvořte handler pomocí generic wrapper

```go
package item

import "awesomeProject2/common"

type Handler struct {
    *common.CRUDHandler[Item]
}

func NewHandler(service *ItemService) *Handler {
    return &Handler{
        CRUDHandler: common.NewCRUDHandler[Item](service),
    }
}
```

**To je vše!** Handler automaticky získá všechny CRUD metody.

### 3. Registrujte routy

```go
itemHandler := item.NewHandler(itemService)

items := router.Group("/items")
{
    items.GET("", itemHandler.GetAll())
    items.GET("/:id", itemHandler.GetByID())
    items.POST("", itemHandler.Create())
    items.PUT("/:id", itemHandler.Update())
    items.DELETE("/:id", itemHandler.Delete())
}
```

## Porovnání s původním kódem

### Před (104 řádky):
```go
func (h *Handler) GetAll(c *gin.Context) {
    items, err := h.service.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, items)
}

func (h *Handler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    item, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, item)
}
// ... další 80+ řádků
```

### Po (18 řádků):
```go
type Handler struct {
    *common.CRUDHandler[Item]
}

func NewHandler(service *ItemService) *Handler {
    return &Handler{
        CRUDHandler: common.NewCRUDHandler[Item](service),
    }
}
```

**Výsledek: 86% méně kódu!**

## Rozšíření funkcionality

Pokud potřebujete vlastní handler metody, jednoduše je přidejte:

```go
type Handler struct {
    *common.CRUDHandler[Item]
    service *ItemService
}

// Vlastní metoda
func (h *Handler) SearchByName() gin.HandlerFunc {
    return func(c *gin.Context) {
        name := c.Query("name")
        items, err := h.service.SearchByName(name)
        if err != nil {
            common.RespondError(c, http.StatusInternalServerError, err)
            return
        }
        common.RespondSuccess(c, http.StatusOK, items)
    }
}
```

## Middleware

Balíček obsahuje také připravené middleware:

```go
// Recovery middleware - zachytává panics
router.Use(common.RecoveryMiddleware())

// Logging middleware - loguje requesty
router.Use(common.LoggingMiddleware())

// CORS middleware - povoluje cross-origin requesty
router.Use(common.CORSMiddleware())
```

## Výhody

- **Méně kódu** - 86% redukce boilerplate kódu
- **Konzistence** - všechny handlery mají stejnou strukturu
- **Type-safe** - díky Go generics
- **Rozšiřitelné** - můžete přidat vlastní metody
- **Testovatelné** - jednodušší mockování
- **Maintainable** - změny v error handlingu na jednom místě

## Utility funkce

```go
// Parsování ID z URL parametru
id, err := common.ParseID(c)

// Úspěšná odpověď
common.RespondSuccess(c, http.StatusOK, data)

// Error odpověď
common.RespondError(c, http.StatusBadRequest, err)
```
