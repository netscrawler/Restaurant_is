package handlers

import (
	"errors"
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

type queryError struct {
	Err  string
	Desc string
	Code int
}

func (e *queryError) Error() string {
	return e.Err + " " + e.Desc
}

func (h *MenuHandler) ListDishes(c *gin.Context) {
	req := &menuv1.ListDishesRequest{}

	var err error

	// Обработка category_id
	if v := c.Query("category_id"); v != "" {
		req.CategoryId, err = processCategoryID(v)
		if err != nil {
			e := &queryError{}
			ok := errors.As(err, &e)
			c.JSON(e.Code, gin.H{
				"error":   e.Err,
				"details": e.Desc,
			})

			return
		}
	}

	// Обработка only_available
	if v := c.Query("only_available"); v != "" {
		req.OnlyAvailable = v == "true" || v == "1"
	}

	// Обработка page с расширенной валидацией
	if v := c.Query("page"); v != "" {
		req.Page, err = processPage(v)
		if err != nil {
			e := &queryError{}
			ok := errors.As(err, &e)
			c.JSON(e.Code, gin.H{
				"error":   e.Err,
				"details": e.Desc,
			})

			return
		}
	}

	// Обработка page_size с расширенной валидацией
	if v := c.Query("page_size"); v != "" {
		req.PageSize, err = processPageSize(v)
		if err != nil {
			e := &queryError{}
			ok := errors.As(err, &e)
			c.JSON(e.Code, gin.H{
				"error":   e.Err,
				"details": e.Desc,
			})

			return
		}
	}

	// Вызов сервиса
	resp, err := h.menuClient.ListDishes(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func processCategoryID(v string) (*int32, error) {
	if len(v) > 10 {
		return nil, &queryError{
			Err:  "error",
			Desc: "category_id is too long",
			Code: http.StatusBadRequest,
		}
	}

	id, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, &queryError{
			Err:  "error",
			Desc: "invalid category_id",
			Code: http.StatusBadRequest,
		}
	}

	if id < 0 {
		return nil, &queryError{
			Err:  "error",
			Desc: "category_id must be positive",
			Code: http.StatusBadRequest,
		}
	}

	val := int32(id)

	return &val, nil
}

func processPage(v string) (int32, error) {
	// Защита от очень длинных строк
	if len(v) > 10 {
		return 0, &queryError{
			Err:  "page value too long",
			Desc: "page number exceeds maximum length",
			Code: http.StatusBadRequest,
		}
	}

	// Проверка на наличие только цифр
	if !isDigitsOnly(v) {
		return 0, &queryError{
			Err:  "invalid page format",
			Desc: "page must contain only digits",
			Code: http.StatusBadRequest,
		}
	}

	page, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, &queryError{
			Err:  "invalid page value",
			Desc: "could not parse page number",
			Code: http.StatusBadRequest,
		}
	}

	// Проверка диапазона
	if page < 1 || page > 10000 { // Реалистичное ограничение для страниц
		return 0, &queryError{
			Err:  "page out of range",
			Desc: "page must be between 1 and 10000",
			Code: http.StatusBadRequest,
		}
	}

	return int32(page), nil
}

func processPageSize(v string) (int32, error) {
	// Защита от очень длинных строк
	if len(v) > 5 {
		return 0, &queryError{
			Err:  "page_size value too long",
			Desc: "page size exceeds maximum length",
			Code: http.StatusBadRequest,
		}
	}

	// Проверка на наличие только цифр
	if !isDigitsOnly(v) {
		return 0, &queryError{
			Err:  "invalid page_size format",
			Desc: "page_size must contain only digits",
			Code: http.StatusBadRequest,
		}
	}

	pageSize, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, &queryError{
			Err:  "invalid page_size value",
			Desc: "could not parse page size",
			Code: http.StatusBadRequest,
		}
	}

	// Проверка диапазона с разумными пределами
	if pageSize < 1 || pageSize > 500 {
		return 0, &queryError{
			Err:  "page_size out of range",
			Desc: "page_size must be between 1 and 500",
			Code: http.StatusBadRequest,
		}
	}

	return int32(pageSize), nil
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

// Вспомогательная функция для проверки, что строка содержит только цифры.
func isDigitsOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
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
