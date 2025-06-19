package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authv1 "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type AuthHandler struct {
	authClient authv1.AuthClient
}

func NewAuthHandler(authClient authv1.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

func (h *AuthHandler) LoginInit(c *gin.Context) {
	var req authv1.LoginClientInitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.LoginClientInit(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LoginConfirm(c *gin.Context) {
	var req authv1.LoginClientConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.LoginClientConfirm(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LoginStaff(c *gin.Context) {
	var req authv1.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.LoginStaff(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) RegisterStaff(c *gin.Context) {
	var req authv1.RegisterStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.RegisterStaff(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LoginYandex(c *gin.Context) {
	var req authv1.OAuthYandexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.LoginYandex(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req authv1.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.Refresh(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Validate(c *gin.Context) {
	var req authv1.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	resp, err := h.authClient.Validate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, resp)
}
