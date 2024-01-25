package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/timohahaa/ewallet/internal/service"
)

type walletRoutes struct {
	walletService service.WalletService
}

func newWalletRoutes(g *echo.Group, ws service.WalletService) {
	r := &walletRoutes{
		walletService: ws,
	}

	g.POST("/wallet", r.CreateWallet)
	g.POST("/wallet/:walletId/send", r.Transfer)
	g.GET("/wallet/:walletId/history", r.TransactionHistory)
	g.GET("/wallet/:walletId", r.Wallet)
}

// POST /api/v1/wallet
func (r *walletRoutes) CreateWallet(c echo.Context) error {
	wallet, err := r.walletService.CreateWallet(c.Request().Context())
	if err != nil {
		newErrorMessage(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

// POST /api/v1/wallet/{walletId}/send
func (r *walletRoutes) Transfer(c echo.Context) error {
	walletId := c.Param("walletId")
	fromWalletId, err := uuid.Parse(walletId)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, "invalid path parametr")
		return err
	}

	var input struct {
		To     uuid.UUID `json:"to"`
		Amount float32   `json:"amount"`
	}
	if err := c.Bind(&input); err != nil {
		newErrorMessage(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	err = r.walletService.Transfer(c.Request().Context(), fromWalletId, input.To, input.Amount)
	if errors.Is(err, service.ErrWalletNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if errors.Is(err, service.ErrTargetWalletNotFound) {
		return c.NoContent(http.StatusBadRequest)
	}
	if errors.Is(err, service.ErrNotEnoughBalance) {
		return c.NoContent(http.StatusBadRequest)
	}
	if err != nil {
		slog.Error("walletRoutes.Transfer - walletService.Transfer", "err", err)
		newErrorMessage(c, http.StatusInternalServerError, "internal server error")
		return nil
	}

	return c.NoContent(http.StatusOK)
}

// GET /api/v1/wallet/{walletId}/history
func (r *walletRoutes) TransactionHistory(c echo.Context) error {
	walletIdStr := c.Param("walletId")
	walletId, err := uuid.Parse(walletIdStr)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, "invalid path parametr")
		return err
	}

	txs, err := r.walletService.TransactionHistory(c.Request().Context(), walletId)
	if errors.Is(err, service.ErrWalletNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		slog.Error("walletRoutes.TransactionHistory - walletService.TransactionHistory", "err", err)
		newErrorMessage(c, http.StatusInternalServerError, "internal server error")
		return nil
	}

	return c.JSON(http.StatusOK, txs)
}

// GET /api/v1/wallet/{walletId}
func (r *walletRoutes) Wallet(c echo.Context) error {
	walletIdStr := c.Param("walletId")
	walletId, err := uuid.Parse(walletIdStr)
	if err != nil {
		newErrorMessage(c, http.StatusBadRequest, "invalid path parametr")
		return err
	}

	wallet, err := r.walletService.WalletStatus(c.Request().Context(), walletId)
	if errors.Is(err, service.ErrWalletNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		slog.Error("walletRoutes.Wallet - walletService.WalletStatus", "err", err)
		newErrorMessage(c, http.StatusInternalServerError, "internal server error")
		return nil
	}

	return c.JSON(http.StatusOK, wallet)
}
