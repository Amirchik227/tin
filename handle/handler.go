package handle

import (
	"fmt"
	"net/http"
	"projectZero/database"
	"projectZero/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) CreateOrder(o storage.Order) int { 
	id := h.storage.Insert(&o)		
	database.InsertOrder(o, id)
	return id
}

func (h *Handler) CreateOrderFromJSON(c *gin.Context) { 
	var order storage.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		fmt.Printf("failed to bind order: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": h.CreateOrder(order), 
	})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	order, err := h.storage.Get(id) 
	if err != nil {
		fmt.Printf("failed to get order %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// func (h *Handler) DeleteOrder(c *gin.Context) {
// 	fmt.Println("DeleteOrder")
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
// 		c.JSON(http.StatusBadRequest, ErrorResponse{
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	h.storage.Delete(id)

// 	c.String(http.StatusOK, "order deleted")
// }

// func (h *Handler) UpdateOrder(c *gin.Context) {
// 	fmt.Println("UpdateOrder")
// id, err := strconv.Atoi(c.Param("id"))
// if err != nil {
// 	fmt.Printf("failed to convert id param to int: %s\n", err.Error())
// 	c.JSON(http.StatusBadRequest, ErrorResponse{
// 		Message: err.Error(),
// 	})
// 	return
// }

// var order Order

// if err := c.BindJSON(&order); err != nil {
// 	fmt.Printf("failed to bind order: %s\n", err.Error())
// 	c.JSON(http.StatusBadRequest, ErrorResponse{
// 		Message: err.Error(),
// 	})
// 	return
// }

// h.storage.Update(id, order)

// c.JSON(http.StatusOK, map[string]interface{}{
// 	"id": order.ID,
// })
// }
