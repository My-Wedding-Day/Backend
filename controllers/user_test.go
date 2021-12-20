package controllers

import (
	"alta-wedding/config"
	"alta-wedding/constants"
	"alta-wedding/lib/database"
	"alta-wedding/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

// Fungsi untuk menginisialisasi koneksi ke database test
func InitEchoTestUser() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

// Struct yang digunakan ketika test request success, dapat menampung banyak data
type UsersResponseSuccess struct {
	Status  string
	Message string
	Data    models.User
}

// type UserResponse struct {
// 	Status  string
// 	Message string
// 	User    models.User
// }

// Struct yang digunakan ketika test request failed
type ResponFailed struct {
	Status  string
	Message string
}

var (
	mock_data_user = models.User{
		Name:     "alterra",
		Email:    "alterra@gmail.com",
		Password: "yourpasswd",
	}
)

var logindata = models.UserLogin{
	Email:    "armuh@gmail.com",
	Password: "yourpass",
}

type ResponSuccessLogin struct {
	Status  string `json:"status" form:"status"`
	Message string `json:"message" form:"message"`
	ID      int    `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Role    string `json:"role" form:"role"`
	Token   string `json:"token" form:"token"`
}

var expass string

// Fungsi untuk memasukkan data user test ke dalam database
func InsertMockDataUserToDB() error {
	expass, _ = database.GeneratehashPassword(mock_data_user.Password)
	mock_data_user.Password = expass
	if err := config.DB.Save(&mock_data_user).Error; err != nil {
		return err
	}
	return nil
}

/*

//------------------------------------------------------
//test get user
func TestGetUserControllers(t *testing.T) {
	testCases := struct {
		name         string
		path         string
		expectStatus int
	}{

		name:         "berhasil",
		path:         "/users/:id",
		expectStatus: 200,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, GetUsersController(context)) {

		var user UserResponse
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectStatus, res.Code)
		assert.Equal(t, "alterra", user.User.Name)

	}
}

//test get user error
func TestGetUserControllersError(t *testing.T) {
	testCases := struct {
		name         string
		path         string
		expectStatus int
	}{

		name:         "User not found",
		path:         "/users/:id",
		expectStatus: http.StatusBadRequest,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(models.User{})
	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, GetUsersController(context)) {

		var user UserResponse
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectStatus, res.Code)
		assert.Equal(t, "Bad Request", user.Message)

	}
}

//------------------------------------------------

//test register user
func TestRegisterUserController(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Success Create User",
		path:       "/users",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()

	body, err := json.Marshal(mock_data_user)
	if err != nil {
		t.Error(t, err, "error")
	}

	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, RegisterUsersController(c)) {
		bodyResponses := rec.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, rec.Code)
		assert.Equal(t, "alterra", user.User.Name)
		assert.Equal(t, "alterra@gmail.com", user.User.Email)
		assert.Equal(t, "Success Create User", user.Message)
	}

}

//test create user error
func TestRegisterUserControllerError(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed to Create User",
		path:       "/users",
		expectCode: http.StatusBadRequest,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(models.User{})

	body, err := json.Marshal(mock_data_user)
	if err != nil {
		t.Error(t, err, "error")
	}

	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, RegisterUsersController(c)) {
		bodyResponses := rec.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, rec.Code)
		assert.Equal(t, "Failed to create user", user.Message)
	}

}

//--------------------------------------------------------
//test update user
func TestUpdateUserControllerSucces(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Success Update User",
		path:       "/user/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	var newdata = models.User{
		Name:     "alta",
		Email:    "alta@gmail.com",
		Password: "qwerty",
	}
	body, err := json.Marshal(newdata)
	if err != nil {
		t.Error(t, err, "error marshal")
	}

	req := httptest.NewRequest(http.MethodPut, "/user/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, UpdateUserController(context)) {
		var response UserResponse
		bodyResponses := res.Body.String()
		err := json.Unmarshal([]byte(bodyResponses), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "Success get user", response.Message)
	}

}

//test update user
func TestUpdateUserControllerParam(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed Update User",
		path:       "/user/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("#")

	if assert.NoError(t, UpdateUserController(context)) {

		var response ResponFailed
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "False Param", response.Message)
	}

}

//test update user error
func TestUpdateUserControllerError(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed Update User",
		path:       "/user/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(&models.User{})
	req := httptest.NewRequest(http.MethodGet, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, UpdateUserController(context)) {

		var response ResponFailed
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "Bad Request", response.Message)
	}

}

//--------------------------------------------------------

//test delete user
func TestDeleteUserController(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Success Delete User",
		path:       "/users/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodDelete, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/users/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")
	if assert.NoError(t, DeleteUserController(contex)) {
		var user ResponFailed
		body := res.Body.String()
		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error unmarshal")
		}
		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "Success Delete User", user.Message)
	}
}

//test delete user error
func TestDeleteUserControllerError(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "Failed to Create User",
		path:       "/users/:id",
		expectCode: http.StatusBadRequest,
	}

	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(&models.User{})
	req := httptest.NewRequest(http.MethodDelete, "/user/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath(testCases.path)
	c.SetParamNames("id")
	c.SetParamValues("1")

	//send data using request body with HTTP Method DELETE
	if assert.NoError(t, DeleteUserController(c)) {
		bodyResponses := res.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "Failed to delete user", user.Message)
	}

}

//-----------------------------------------
//tes login error
func TestLoginUsersControllersError(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	config.DB.Migrator().DropTable(models.User{})
	req := httptest.NewRequest(http.MethodPost, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)

	if assert.NoError(t, LoginUsersController(context)) {

		var user UserResponse
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		// assert.Equal(t, testCases.expectStatus, res.Code)
		assert.Equal(t, "Login failed", user.Message)

	}
}

*/

//test login success
func TestLoginGetUserControllers(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockDataUserToDB()
	body, error := json.Marshal(logindata)
	if error != nil {
		t.Error(t, error, "error marshal")
	}
	// send data using request body with HTTP method POST
	req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/login/users")
	middleware.JWT([]byte(constants.SECRET_JWT))(LoginUsersControllerTest())(context)

	if assert.NoError(t, LoginUsersController(context)) {
		res_body := res.Body.String()
		var User ResponSuccessLogin
		err := json.Unmarshal([]byte(res_body), &User)
		if err != nil {
			assert.Error(t, err, "error marshal")
		}

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, 1, User.ID)
		assert.Equal(t, "success", User.Status)
		assert.Equal(t, "login success", User.Message)
		assert.Equal(t, "Arif Muhammad", User.Name)
		assert.Equal(t, "user", User.Role)

	}
}

//test login error
func TestLoginUserFailed(t *testing.T) {
	e := InitEchoTestUser()
	InsertMockDataUserToDB()

	t.Run("TestLoginUser_InvalidInput", func(t *testing.T) {
		logininfo, err := json.Marshal(models.UserLogin{Email: "fian@gmail.com", Password: "admins"})
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(logininfo))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/users")
		if assert.NoError(t, LoginUsersController(contex)) {
			bodyResponses := rec.Body.String()
			var User ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &User)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", User.Status)
			assert.Equal(t, "invalid email or password", User.Message)
		}
	})
	t.Run("TestLoginUser_ErrorDB", func(t *testing.T) {
		datalogin, err := json.Marshal(logindata)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		req := httptest.NewRequest(http.MethodPost, "/login/users", bytes.NewBuffer(datalogin))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/users")
		config.DB.Migrator().DropTable(&models.User{})
		if assert.NoError(t, LoginUsersController(contex)) {
			bodyResponses := rec.Body.String()
			var User ResponFailed
			err := json.Unmarshal([]byte(bodyResponses), &User)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "failed", User.Status)
			assert.Equal(t, "internal server error", User.Message)
		}
	})
}
