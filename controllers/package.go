package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

// Controller untuk memasukkan package baru
func InsertPackageController(c echo.Context) error {
	// Mendapatkan data package baru dari client
	input := models.Package{}
	c.Bind(&input)
	duplicate, _ := database.GetPackageByName(input.PackageName)
	if duplicate > 0 {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("package name was use, try input another package name"))
	}
	organizer_id := middlewares.ExtractTokenUserId(c)
	input.Organizer_ID = organizer_id
	// Menyimpan data barang baru menggunakan fungsi InsertPackage
	data, e := database.InsertPackage(input)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("failed to input package"))
	}

	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("urlphoto")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	fileExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "bmp": true}
	ext := strings.Split(uploaded_file.Filename, ".")
	extension := ext[len(ext)-1]
	if !fileExtensions[extension] {
		return c.JSON(http.StatusBadRequest, responses.StatusFailedDataPhoto("invalid type"))
	}

	t := time.Now()
	formatted := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	packageName := strings.ReplaceAll(input.PackageName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("%v-%v.%v", packageName, formatted, extension)
	sw := storageClient.Bucket(bucket).Object(uploaded_file.Filename).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	if err := sw.Close(); err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	u, err := url.Parse("https://storage.googleapis.com/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}

	// Insert URL
	urlPhoto := fmt.Sprintf("%v", u)
	foto := models.Photo{
		Package_ID: data.ID,
		Photo_Name: packageName,
		UrlPhoto:   urlPhoto,
	}
	_, tx := database.InsertPhoto(foto)
	if tx != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}

	return c.JSON(http.StatusCreated, responses.StatusSuccess("success to input package"))
}

// Controller untuk mendapatkan seluruh data Packages
func GetAllPackageController(c echo.Context) error {
	// Mendapatkan data satu buku menggunakan fungsi GetPackages
	paket, e := database.GetPackages()
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch packages"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all packages", paket))
}

// // Controller untuk mendapatkan seluruh data Packages by token
// func GetAllPackageByTokenController(c echo.Context) error {
// 	idOrganizer := middlewares.ExtractTokenOrganizerId(c)
// 	// Mendapatkan data satu buku menggunakan fungsi GetPackages
// 	paket, e := database.GetPackagesByToken(idOrganizer)
// 	if e != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch packages"))
// 	}
// 	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all packages by token", paket))
// }

// Controller untuk mendapatkan seluruh data Packages by ID
func GetPackageByIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	// Mendapatkan data satu buku menggunakan fungsi GetPackages
	paket, e := database.GetPackagesByID(id)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to fetch packages"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get all packages by ID", paket))
}
