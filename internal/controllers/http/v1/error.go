package v1

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrTargetWalletNotFound = errors.New("target wallet not found")
	ErrNotEnoughBalance     = errors.New("not enough balance")
)

func newErrorMessage(c echo.Context, statusCode int, message string) {
	httpErr := echo.NewHTTPError(statusCode, message)
	_ = c.JSON(statusCode, httpErr)
}
