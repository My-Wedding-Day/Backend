package controllers

// import (
// 	"alta-wedding/lib/database"
// 	responses "alta-wedding/lib/responses"
// 	"alta-wedding/middlewares"
// 	"alta-wedding/models"
// 	"net/http"
// 	"time"

// 	"github.com/labstack/echo/v4"
// 	"golang.org/x/tools/go/packages"
// )

// Controller untuk memasukkan barang baru ke Reservation
// func CreateReservationController(c echo.Context) error {
// 	Reservation := models.Reservation{}
// 	c.Bind(&Reservation)
// 	logged := middlewares.ExtractTokenUserId(c)

// 	input, _ := database.GetPackagesByID(Reservation.Package_ID)

// 	reservation, err := database.CreateReservation(input, int(input.Package_ID))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed())
// 	}

// 	return c.JSON(http.StatusCreated, responses.ReservationSuccess())

// }
