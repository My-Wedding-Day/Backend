package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	input := models.PostRequestBodyPackage{}
	c.Bind(&input)
	duplicate, _ := database.GetPackageByName(input.PackageName)
	if duplicate > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("package name was user, try input another package name"))
	}

	paket := models.Package{
		PackageName: input.PackageName,
		Price:       input.Price,
		Pax:         input.Pax,
		PackageDesc: input.PackageDesc,
	}

	// Menyimpan data barang baru menggunakan fungsi InsertPackage
	data, e := database.InsertPackage(paket)
	if e != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("failed to input package"))
	}

	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("url")
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
		Url:        urlPhoto,
	}
	_, tx := database.InsertPhoto(foto)
	if tx != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}

	return c.JSON(http.StatusOK, responses.StatusSuccess("success to input package"))
}
