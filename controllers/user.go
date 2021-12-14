package controllers

import (
	"alta-wedding/lib/database"
	responses "alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

//register user
func RegisterUsersController(c echo.Context) error {
	var user models.User
	c.Bind(&user)
	duplicate, _ := database.GetUserByEmail(user.Email)
	if duplicate > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Email was used, try another email"))
	}

	Password, _ := database.GeneratehashPassword(user.Password)
	user.Password = Password
	user.Role = "User"
	_, err := database.RegisterUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create new user"))
}

//login users
func LoginUsersController(c echo.Context) error {
	login := models.UserLogin{}
	c.Bind(&login)
	users, err := database.LoginUsers(&login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid email"))
	}
	if users == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid password"))
	}
	token, err := middlewares.CreateToken(int(users.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("can not generate token"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccessLogin("login success", users.ID, token, users.Name, users.Role))
}

//get user by id
func GetUsersController(c echo.Context) error {
	loginuser := middlewares.ExtractTokenUserId(c)
	datauser, e := database.GetUser(loginuser)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	respon := models.Profile{
		ID:    datauser.ID,
		Name:  datauser.Name,
		Email: datauser.Email,
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get user", respon))
}

//update user by id
func UpdateUserController(c echo.Context) error {
	var user models.User
	c.Bind(&user)
	loginuser := middlewares.ExtractTokenUserId(c)
	_, e := database.UpdateUser(loginuser, user)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal service error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccess("success update user"))
}
