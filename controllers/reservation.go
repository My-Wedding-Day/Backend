package controllers

import (
	"alta-wedding/lib/database"
	responses "alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller untuk memasukkan barang baru ke Reservation
func CreateReservationController(c echo.Context) error {
	Reservation := models.Reservation{}
	c.Bind(&Reservation)
	// EXTRACT TOKEN LOGIN
	logged := middlewares.ExtractTokenUserId(c)
	// GET PACKAGE ID
	input, _ := database.GetPackagesByID(Reservation.Package_ID)
	// WRONG INPUT
	if input == nil {
		return c.JSON(http.StatusBadRequest, responses.ReservationFailed())
	}
	// TRANSFER DATA FROM TOKEN
	Reservation.User_ID = logged
	// CREATE RESERVATION
	respon, err := database.CreateReservation(&Reservation)
	// DATABASE OR SERVER INTERNAL ERROR
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	// DATABASE ALREADY RESERVATION
	if respon == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("already reserve"))
	}
	// RESERVATION SUCCESS
	return c.JSON(http.StatusCreated, responses.ReservationSuccess())
}

func GetReservationController(c echo.Context) error {
	logged := middlewares.ExtractTokenUserId(c)
	input, e := database.GetReservation(logged)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success", input))
}
