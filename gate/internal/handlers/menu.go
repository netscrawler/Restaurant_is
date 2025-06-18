package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	menuv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/menu"
)

type MenuHandler struct {
	menuClient menuv1.MenuServiceClient
}

func NewMenuHandler(menuClient menuv1.MenuServiceClient) *MenuHandler {
	return &MenuHandler{
		menuClient: menuClient,
	}
}

func (h *MenuHandler) ListDishes(c *gin.Context) {
	req := &menuv1.ListDishesRequest{}

	if v := c.Query("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			val := int32(id)
			req.CategoryId = &val
		}
	}
	if v := c.Query("only_available"); v != "" {
		if v == "true" || v == "1" {
			req.OnlyAvailable = true
		}
	}
	if v := c.Query("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil {
			req.Page = int32(page)
		}
	}
	if v := c.Query("page_size"); v != "" {
		if pageSize, err := strconv.Atoi(v); err == nil {
			req.PageSize = int32(pageSize)
		}
	}

	resp, err := h.menuClient.ListDishes(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) GetDish(c *gin.Context) {
	id := c.Param("id")

	req := &menuv1.GetDishRequest{
		DishId: &menuv1.UUID{Value: id},
	}

	resp, err := h.menuClient.GetDish(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) CreateDish(c *gin.Context) {
	var req menuv1.DishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.menuClient.CreateDish(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MenuHandler) UpdateDish(c *gin.Context) {
	id := c.Param("id")
	var req menuv1.UpdateDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = &menuv1.UUID{Value: id}

	resp, err := h.menuClient.UpdateDish(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) DeleteDish(c *gin.Context) {
	id := c.Param("id")

	req := &menuv1.GetDishRequest{
		DishId: &menuv1.UUID{Value: id},
	}

	_, err := h.menuClient.GetDish(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
