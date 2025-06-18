package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orderv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/order"
)

type OrderHandler struct {
	orderClient orderv1.OrderServiceClient
}

func NewOrderHandler(orderClient orderv1.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderv1.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем user_id из контекста
	userID, exists := c.Get("user_id")
	if exists {
		req.UserId = &orderv1.UUID{Value: userID.(string)}
	}

	resp, err := h.orderClient.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	req := &orderv1.GetOrderRequest{
		OrderId: &orderv1.UUID{Value: id},
	}

	resp, err := h.orderClient.GetOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	req := &orderv1.ListOrdersRequest{}

	// Получаем user_id из контекста для фильтрации
	userID, exists := c.Get("user_id")
	if exists {
		req.UserId = &orderv1.UUID{Value: userID.(string)}
	}

	resp, err := h.orderClient.ListOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req orderv1.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OrderId = &orderv1.UUID{Value: id}

	_, err := h.orderClient.UpdateOrderStatus(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
