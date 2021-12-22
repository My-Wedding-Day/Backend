package controllers

import (
	"alta-wedding/config"
	"alta-wedding/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func InitEchoTestAPIPackage() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

var (
	mock_data_package_tanpa_foto = models.Package{
		Organizer_ID: 1,
		PackageName:  "Coba2",
		Price:        15000000,
		Pax:          100,
		PackageDesc:  "Package Desc",
	}
)

func InsertMockDataPackageTanpaFotoToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_package_tanpa_foto).Error; err != nil {
		return err
	}
	return nil
}

func TestGetPackageByIDSuccess(t *testing.T) {
	e := InitEchoTestAPI()
	InsertMockDataPackageTanpaFotoToDB()
	req := httptest.NewRequest(http.MethodGet, "/package/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	contex := e.NewContext(req, res)
	contex.SetPath("/package/:id")
	contex.SetParamNames("id")
	contex.SetParamValues("1")

	if assert.NoError(t, GetPackageByIDController(contex)) {
		var paket PackageResponSuccess
		body := res.Body.String()
		json.Unmarshal([]byte(body), &paket)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "success get packages by ID", paket.Message)
		assert.Equal(t, "success", paket.Status)
		assert.Equal(t, "Coba2", paket.Data[0].PackageName)
		assert.Equal(t, 100, paket.Data[0].Pax)

	}
}
