package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/timohahaa/ewallet/internal/service"
)

func NewRouter(walletService service.WalletService) *echo.Echo {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	v1 := e.Group("/api/v1")
	{
		newWalletRoutes(v1, walletService)
	}

	return e
}
