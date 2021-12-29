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

func InitEchoTestAPIPackage() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

type PackageSingleResponSuccess struct {
	Status  string
	Message string
	Data    models.Package
}

type PackageManyResponSuccess struct {
	Status  string
	Message string
	Data    []database.GetPackageAllStruct
}

var (
	mock_data_package_tanpa_foto = models.Package{
		Organizer_ID: 1,
		PackageName:  "Coba",
		Price:        15000000,
		Pax:          100,
		PackageDesc:  "Package Desc",
	}
	mock_data_package_tanpa_foto_update = models.Package{
		Organizer_ID: 1,
		PackageName:  "Coba1",
		Price:        1,
		Pax:          1000,
		PackageDesc:  "Package Desc Baru",
	}
)

// // Fungsi untuk melakukan login dan ekstraksi token JWT
// func UsingJWTCart() (string, error) {
// 	// Melakukan login data user test
// 	InsertMockDataOrganizerToDB()
// 	var user models.Organizer
// 	tx := config.DB.Where("email = ? AND password = ?", mock_data_organizer.Email, mock_data_organizer.Password).First(&user)
// 	if tx.Error != nil {
// 		return "", tx.Error
// 	}
// 	// Mengektraksi token data user test
// 	token, err := middlewares.CreateToken(int(user.ID))
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

func InsertMockDataPackageTanpaFotoToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_package_tanpa_foto).Error; err != nil {
		return err
	}
	return nil
}

func InsertMockDataPackageTanpaFotoUpdateToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_package_tanpa_foto_update).Error; err != nil {
		return err
	}
	return nil
}

func TestGetPackageByIDSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	req := httptest.NewRequest(http.MethodGet, "/package/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")

	if assert.NoError(t, GetPackageByIDController(contex)) {
		var paket PackageSingleResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "success get packages by ID", paket.Message)
		assert.Equal(t, "success", paket.Status)
		// assert.Equal(t, "Coba", paket.Data.PackageName)
		// assert.Equal(t, 100, paket.Data.Pax)
	}
}

func TestGetPackageByIDFail(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	req := httptest.NewRequest(http.MethodGet, "/package/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("#")

	if assert.NoError(t, GetPackageByIDController(contex)) {
		var paket PackageSingleResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "false param", paket.Message)
		assert.Equal(t, "failed", paket.Status)
		// assert.Equal(t, "Coba", paket.Data.PackageName)
		// assert.Equal(t, 100, paket.Data.Pax)
	}
}

func TestGetPackageByIDFailFetch(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	config.DB.Migrator().DropTable(&models.Package{})
	req := httptest.NewRequest(http.MethodGet, "/package/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")

	if assert.NoError(t, GetPackageByIDController(contex)) {
		var paket PackageSingleResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "failed to fetch packages", paket.Message)
		assert.Equal(t, "failed", paket.Status)
		// assert.Equal(t, "Coba", paket.Data.PackageName)
		// assert.Equal(t, 100, paket.Data.Pax)
	}
}

func TestGetPackageALLSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	req := httptest.NewRequest(http.MethodGet, "/package", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package")
	// contex.SetParamNames("id")
	// contex.SetParamValues("1")

	if assert.NoError(t, GetAllPackageController(contex)) {
		var paket PackageManyResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "success get all packages", paket.Message)
		assert.Equal(t, "success", paket.Status)
		// assert.Equal(t, "Coba", paket.Data[0].PackageName)
		// assert.Equal(t, 100, paket.Data[0].Pax)
	}
}

func TestGetPackageALLFailed(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	config.DB.Migrator().DropTable(&models.Package{})
	req := httptest.NewRequest(http.MethodGet, "/package", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package")
	// contex.SetParamNames("id")
	// contex.SetParamValues("1")

	if assert.NoError(t, GetAllPackageController(contex)) {
		var paket PackageManyResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "failed to fetch packages", paket.Message)
		assert.Equal(t, "failed", paket.Status)
		// assert.Equal(t, "Coba", paket.Data[0].PackageName)
		// assert.Equal(t, 100, paket.Data[0].Pax)
	}
}

func TestUpdatePackageSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	// Mendapatkan data update package
	body, err := json.Marshal(mock_data_package_tanpa_foto_update)
	if err != nil {
		t.Error(t, err, "error")
	}
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodPut, "/package/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/package/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdatePackageControllerTest())(context)

	var paket PackageSingleResponSuccess
	bodyReq := res.Body.String()
	json.Unmarshal([]byte(bodyReq), &paket)
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, "success edit data", paket.Message)
}
