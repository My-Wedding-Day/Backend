package routers

import (
	"alta-wedding/constants"
	"alta-wedding/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	//CORS
	e.Use(echoMid.CORSWithConfig(echoMid.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete},
	}))

	// Midleware Auth
	r := e.Group("")
	r.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	// ------------------------------------------------------------------
	// LOGIN & REGISTER ORGANIZER
	// ------------------------------------------------------------------
	e.POST("/register/organizer", controllers.CreateOrganizerController)
	e.POST("/login/organizer", controllers.LoginOrganizerController)
	r.PUT("/organizer/profile", controllers.UpdateOrganizerController)
	r.PUT("/organizer/profile/photo", controllers.UpdatePhotoOrganizerController)
	return e
}
