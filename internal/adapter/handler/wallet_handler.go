package handler

import (
	"context"
	"hexagonal_intro/internal/adapter/handler/dto"
	"hexagonal_intro/internal/core/model"
	"hexagonal_intro/internal/core/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req dto.CreateWalletRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	wallet, err := h.svc.CreateWallet(context.Background(), req.UserID)
	if err != nil {
		c.JSON(http.StatusConflict, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toWalletResponse(wallet))
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	id := c.Param("id")

	wallet, err := h.svc.GetWallet(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, toWalletResponse(wallet))
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	id := c.Param("id")
	var req dto.TransactionRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}
	wallet, err := h.svc.Deposit(context.Background(), id, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, toWalletResponse(wallet))
}

func (h *WalletHandler) WithDraw(c *gin.Context) {
	id := c.Param("id")

	var req dto.TransactionRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	wallet, err := h.svc.WithDraw(context.Background(), id, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, toWalletResponse(wallet))
}

func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteWallet(context.Background(), id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponseDto{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func toWalletResponse(w *model.Wallet) dto.WalletResponseDto {
	return dto.WalletResponseDto{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   w.Balance,
		CreatedAt: w.CreateAt.Format("2006-01-02T15:04:05z07:00"), //time.RFC3339
		UpdatedAt: w.CreateAt.Format("2006-01-02T15:04:05z07:00"), //time.RFC3339
	}
}
