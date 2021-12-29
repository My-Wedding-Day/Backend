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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
	mock_data_foto = models.Photo{
		Package_ID: 1,
		Photo_Name: "Coba",
		UrlPhoto:   "./cydra.png",
	}
	mock_data_organizer2 = models.Organizer{
		WoName:      "coba2wedd",
		Email:       "coba@coba.coba",
		Password:    "yourpass",
		PhoneNumber: "081232323",
		City:        "Makassar",
		Address:     "Jl. Kertajaya",
	}
)

var logininfo2 = models.LoginRequestBody{
	Email:    "coba@coba.coba",
	Password: "yourpass",
}

func InsertMockDataPackageTanpaFotoToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_package_tanpa_foto).Error; err != nil {
		return err
	}
	return nil
}

func InsertMockDataOrganizer2ToDB() error {
	xpassOrganizer, _ = database.GeneratehashPassword(mock_data_organizer2.Password)
	mock_data_organizer2.Password = xpassOrganizer
	var err error
	if err = config.DB.Save(&mock_data_organizer2).Error; err != nil {
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

func InsertMockDataFotoToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_foto).Error; err != nil {
		return err
	}
	return nil
}

// func InsertMockDataFoto2ToDB() error {
// 	var err error
// 	if err = config.DB.Save(&mock_data_foto2).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

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
	InsertMockDataFotoToDB()
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

func TestUpdatePackageFalseParam(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
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
	context.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdatePackageControllerTest())(context)

	var paket PackageSingleResponSuccess
	bodyReq := res.Body.String()
	json.Unmarshal([]byte(bodyReq), &paket)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "false param", paket.Message)
}

func TestUpdatePackageFalseID(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataOrganizer2ToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
	// Mendapatkan data update package
	body, err := json.Marshal(mock_data_package_tanpa_foto_update)
	if err != nil {
		t.Error(t, err, "error")
	}
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo2.Email, xpassOrganizer).First(&organizerDetail)
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
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Equal(t, "Unauthorized Access", paket.Message)
}

func TestDeletePackageByIDSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodDelete, "/package/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeletePackageControllerTest())(contex)

	var paket PackageSingleResponSuccess
	body := res.Body.String()
	json.Unmarshal([]byte(body), &paket)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "success deleted package", paket.Message)
	assert.Equal(t, "success", paket.Status)
	// assert.Equal(t, "Coba", paket.Data.PackageName)
	// assert.Equal(t, 100, paket.Data.Pax)
}

func TestDeletePackageByIDFalseParam(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodDelete, "/package/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("#")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeletePackageControllerTest())(contex)

	var paket PackageSingleResponSuccess
	body := res.Body.String()
	json.Unmarshal([]byte(body), &paket)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "false param", paket.Message)
	assert.Equal(t, "failed", paket.Status)
	// assert.Equal(t, "Coba", paket.Data.PackageName)
	// assert.Equal(t, 100, paket.Data.Pax)
}

func TestDeletePackageFalseID(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataOrganizer2ToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()

	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo2.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodPut, "/package/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/package/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(DeletePackageControllerTest())(context)

	var paket PackageSingleResponSuccess
	bodyReq := res.Body.String()
	json.Unmarshal([]byte(bodyReq), &paket)
	assert.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Equal(t, "Unauthorized Access", paket.Message)
}

func TestDeletePackageFailedFetch(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()

	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}
	req := httptest.NewRequest(http.MethodPut, "/package/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/package/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")
	config.DB.Migrator().DropTable(&models.Package{})
	config.DB.Migrator().DropTable(&models.Photo{})
	config.DB.Migrator().DropTable(&models.Organizer{})
	middleware.JWT([]byte(constants.SECRET_JWT))(DeletePackageControllerTest())(context)

	var paket PackageSingleResponSuccess
	bodyReq := res.Body.String()
	json.Unmarshal([]byte(bodyReq), &paket)
	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "failed to fetch package", paket.Message)
}

func TestUpdatePhotoPackageControllerSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}

	path := mock_data_foto.UrlPhoto

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("urlphoto", path)
	assert.NoError(t, err)
	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPut, "/package/photo/:id", body)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/package/photo/:id")
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(UpdatePhotoPackageControllerTest())(context)

	var Photo ResponSuccess
	bodyReq := rec.Body.String()
	json.Unmarshal([]byte(bodyReq), &Photo)
	assert.Equal(t, "success upload photo", Photo.Message)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestInsertPackageControllerSuccess(t *testing.T) {
	e := InitEchoTestAPIPackage()
	InsertMockDataOrganizerToDB()
	InsertMockDataPackageTanpaFotoToDB()
	InsertMockDataFotoToDB()
	var organizerDetail models.Organizer
	tx := config.DB.Where("email=? AND password=?", logininfo.Email, xpassOrganizer).First(&organizerDetail)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(organizerDetail.ID))
	if err != nil {
		t.Error("error create token")
	}

	path := mock_data_foto.UrlPhoto

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("urlphoto", path)
	assert.NoError(t, err)
	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/package", body)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/package")
	middleware.JWT([]byte(constants.SECRET_JWT))(InsertPackageControllerTest())(context)

	var Photo ResponSuccess
	bodyReq := rec.Body.String()
	json.Unmarshal([]byte(bodyReq), &Photo)
	assert.Equal(t, "success to input package", Photo.Message)
	assert.Equal(t, http.StatusCreated, rec.Code)

}
