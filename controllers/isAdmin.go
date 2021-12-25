package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetInvoiceController(c echo.Context) error {
	invoiceadmin := middlewares.ExtractTokenUserId(c)
	datauser, e := database.GetUser(invoiceadmin)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if datauser.Role != "admin" {
		return c.JSON(http.StatusUnauthorized, responses.StatusUnauthorized())
	}
	payment, err := database.GetInvoiceAdmin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get admin", payment))
}
