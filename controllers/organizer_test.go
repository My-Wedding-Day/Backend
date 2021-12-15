package controllers

import (
	"alta-wedding/config"
	"alta-wedding/constants"
	"alta-wedding/lib/database"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func InitEchoTestAPI() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

type OrganizerResponSuccess struct {
	Status  string
	Message string
	Data    models.Organizer
}

type ResponseFailed struct {
	Status  string
	Message string
}

type PostErrorBind struct {
	WoName int
}

var logininfo = models.LoginRequestBody{
	Email:    "fian@gmail.com",
	Password: "yourpass",
}

var (
	mock_data_organizer = models.Organizer{
		WoName:      "Makassar Wedding",
		Email:       "fian@gmail.com",
		Password:    "yourpass",
		PhoneNumber: "081242323123",
		City:        "Makassar",
		Address:     "Jl. Penghibur",
	}
)

type LoginResponSuccess struct {
	Status  string `json:"status" form:"status"`
	Message string `json:"message" form:"message"`
	ID      int    `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Role    string `json:"role" form:"role"`
	Token   string `json:"token" form:"token"`
}

var xpass string

func InsertMockDataOrganizerToDB() error {
	xpass, _ = database.GeneratehashPassword(mock_data_organizer.Password)
	mock_data_organizer.Password = xpass
	var err error
	if err = config.DB.Save(&mock_data_organizer).Error; err != nil {
		return err
	}
	return nil
}

func TestLoginOrganizerSuccess(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()
	datalogin, err := json.Marshal(logininfo)
	if err != nil {
		t.Error(t, err, "error marshal")
	}
	// send data using request body with HTTP method POST
	req := httptest.NewRequest(http.MethodPost, "/login/organizer", bytes.NewBuffer(datalogin))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	contex := e.NewContext(req, rec)
	contex.SetPath("/login/organizer")
	if assert.NoError(t, LoginOrganizerController(contex)) {
		bodyResponses := rec.Body.String()
		var organizer LoginResponSuccess
		err := json.Unmarshal([]byte(bodyResponses), &organizer)
		if err != nil {
			assert.Error(t, err, "error marshal")
		}
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, 1, organizer.ID)
		assert.Equal(t, "success", organizer.Status)
		assert.Equal(t, "login success", organizer.Message)
		assert.Equal(t, "Makassar Wedding", organizer.Name)
		assert.Equal(t, "organizer", organizer.Role)
	}
}

func TestLoginOrganizerFailed(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()

	t.Run("TestLoginOrganizer_InvalidInput", func(t *testing.T) {
		logininfo, err := json.Marshal(models.LoginRequestBody{Email: "fian@gmail.com", Password: "admins"})
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/login/organizer", bytes.NewBuffer(logininfo))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/organizer")
		if assert.NoError(t, LoginOrganizerController(contex)) {
			bodyResponses := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(bodyResponses), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "invalid email or password", organizer.Message)
		}
	})
	t.Run("TestLoginOrganizer_ErrorDB", func(t *testing.T) {
		datalogin, err := json.Marshal(logininfo)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		req := httptest.NewRequest(http.MethodPost, "/login/organizer", bytes.NewBuffer(datalogin))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		contex.SetPath("/login/organizer")
		config.DB.Migrator().DropTable(&models.Organizer{})
		if assert.NoError(t, LoginOrganizerController(contex)) {
			bodyResponses := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(bodyResponses), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "internal server error", organizer.Message)
		}
	})
}

func TestGetProfileByIDOrganizerSuccess(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()
	req := httptest.NewRequest(http.MethodGet, "/organizer/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/organizer/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")

	if assert.NoError(t, GetProileOrganizerbyIDController(contex)) {
		var organizer OrganizerResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &organizer)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "success get organizer", organizer.Message)
		assert.Equal(t, "success", organizer.Status)
		assert.Equal(t, "Makassar Wedding", organizer.Data.WoName)
		assert.Equal(t, "fian@gmail.com", organizer.Data.Email)
		assert.Equal(t, "Makassar", organizer.Data.City)

	}
}
func TestGetProfileByIDOrganizerFailed(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()
	t.Run("TestGetByIDProfile_FalseParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/organizer/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		contex := e.NewContext(req, res)
		contex.SetPath("/organizer/:id")
		contex.SetParamNames("id")
		contex.SetParamValues("#")

		if assert.NoError(t, GetProileOrganizerbyIDController(contex)) {
			var organizer OrganizerResponSuccess
			body := res.Body.String()
			json.Unmarshal([]byte(body), &organizer)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Equal(t, "false param", organizer.Message)
		}
	})
	t.Run("TestGetByIDProfile_ErrorDB", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/organizer/:id", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		contex := e.NewContext(req, res)
		contex.SetPath("/organizer/:id")
		contex.SetParamNames("id")
		contex.SetParamValues("1")
		config.DB.Migrator().DropTable(&models.Organizer{})
		if assert.NoError(t, GetProileOrganizerbyIDController(contex)) {
			var organizer OrganizerResponSuccess
			body := res.Body.String()
			json.Unmarshal([]byte(body), &organizer)
			assert.Equal(t, http.StatusInternalServerError, res.Code)
			assert.Equal(t, "internal server error", organizer.Message)
		}
	})
}

func TestGetProfileOrganizerSuccess(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpass).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/users/profile")
	middleware.JWT([]byte(constants.SECRET_JWT))(GetProfileOrganizerControllerTest())(context)

	var organizer OrganizerResponSuccess
	body := res.Body.String()
	json.Unmarshal([]byte(body), &organizer)
	t.Run("GET/organizer/profile", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "success get organizer", organizer.Message)
		assert.Equal(t, "success", organizer.Status)
		assert.Equal(t, "Makassar Wedding", organizer.Data.WoName)
		assert.Equal(t, "fian@gmail.com", organizer.Data.Email)
		assert.Equal(t, "Makassar", organizer.Data.City)
	})
}

func TestGetProfileOrganizerFailed(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataOrganizerToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpass).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/users/profile")
	// drop table from database
	config.DB.Migrator().DropTable(&models.Organizer{})
	middleware.JWT([]byte(constants.SECRET_JWT))(GetProfileOrganizerControllerTest())(context)
	var organizer OrganizerResponSuccess
	body := res.Body.String()
	json.Unmarshal([]byte(body), &organizer)
	t.Run("GET/organizer/profile", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, "failed", organizer.Status)
		assert.Equal(t, "internal server error", organizer.Message)
	})
}

func TestRegisterOrganizerSuccess(t *testing.T) {
	e := InitEchoTestAPI()
	body, err := json.Marshal(mock_data_organizer)
	if err != nil {
		t.Error(t, err, "error marshal")
	}
	// send data using request body with HTTP method POST
	req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	contex := e.NewContext(req, rec)
	if assert.NoError(t, CreateOrganizerController(contex)) {
		body := rec.Body.String()
		var organizer OrganizerResponSuccess
		err := json.Unmarshal([]byte(body), &organizer)
		if err != nil {
			assert.Error(t, err, "error marshal")
		}
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "success", organizer.Status)
		assert.Equal(t, "success create new organizer", organizer.Message)
	}
}

func TestRegisterOrganizerFailed(t *testing.T) {
	e := InitEchoTestAPI()
	t.Run("TestRegister_ErrorBind", func(t *testing.T) {
		body, err := json.Marshal(PostErrorBind{})
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "bad request", organizer.Message)
		}
	})
	t.Run("TestRegister_InputEmpty", func(t *testing.T) {
		type Login struct {
			Name     string
			Password string
		}
		var empty Login
		body, err := json.Marshal(empty)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "input data cannot be empty", organizer.Message)
		}
	})
	t.Run("TestRegister_WONameLess", func(t *testing.T) {
		mock_data_organizer.WoName = "Fian"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "business name cannot less than 5 characters", organizer.Message)
		}
	})
	t.Run("TestRegister_EmailWasUsed", func(t *testing.T) {
		InsertMockDataOrganizerToDB()
		mock_data_organizer.WoName = "Makassar wedding"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "email was used, try another one", organizer.Message)
		}
	})
	t.Run("TestRegister_WoNameUsed", func(t *testing.T) {
		InsertMockDataOrganizerToDB()
		mock_data_organizer.WoName = "Makassar wedding"
		mock_data_organizer.Email = "alfian1@gmail.com"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "business name was used, try another one", organizer.Message)
		}
	})
	t.Run("TestRegister_InvalidFormatName", func(t *testing.T) {
		mock_data_organizer.WoName = "     Makassar wedding"
		mock_data_organizer.Email = "alfian1123@gmail.com"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "invalid format name", organizer.Message)
		}
	})
	t.Run("TestRegister_InvalidFormatEmail", func(t *testing.T) {
		mock_data_organizer.WoName = "Makassar wedding Exlusive"
		mock_data_organizer.Email = "#alfian@gmail.com"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "email must contain email format", organizer.Message)
		}
	})
	t.Run("TestRegister_LengthCharacter", func(t *testing.T) {
		mock_data_organizer.Password = "inipass"
		mock_data_organizer.Email = "alfian@email.com"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "password or phone number cannot less than 8 characters", organizer.Message)
		}
	})
	t.Run("TestRegister_DuplicatePhone", func(t *testing.T) {
		mock_data_organizer.Password = "inipassword"
		mock_data_organizer.Email = "alfian@email.com"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "phone number was used, try another one", organizer.Message)
		}
	})
	t.Run("TestRegister_ErrorDB", func(t *testing.T) {
		config.DB.Migrator().DropTable(&models.Organizer{})
		mock_data_organizer.Password = "inipassword"
		mock_data_organizer.Email = "alfianfian@gmail.com"
		mock_data_organizer.PhoneNumber = "12345678910"
		body, err := json.Marshal(mock_data_organizer)
		if err != nil {
			t.Error(t, err, "error marshal")
		}
		// send data using request body with HTTP method POST
		req := httptest.NewRequest(http.MethodPost, "/register/organizer", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		contex := e.NewContext(req, rec)
		if assert.NoError(t, CreateOrganizerController(contex)) {
			body := rec.Body.String()
			var organizer ResponseFailed
			err := json.Unmarshal([]byte(body), &organizer)
			if err != nil {
				assert.Error(t, err, "error marshal")
			}
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "failed", organizer.Status)
			assert.Equal(t, "internal server error", organizer.Message)
		}
	})

}
