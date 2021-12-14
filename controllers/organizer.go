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
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

// Register Organizer Function
func CreateOrganizerController(c echo.Context) error {
	organizer := models.Organizer{}
	// Bind all data from JSON
	if err := c.Bind(&organizer); err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	// Check Organizer is Exist
	row, _ := database.FindOrganizerByEmail(organizer.Email)
	if row != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email was used, try another email"))
	}
	// hash password bcrypt
	password, _ := database.GeneratehashPassword(organizer.Password)
	organizer.Password = password // replace old password to bcrypt password
	// Insert ALL data to Database
	_, e := database.InsertOrganizer(organizer)
	if e != nil {
		// Respon Failed
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	// Respon Success
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success create new organizer"))
}

// Login Organizer Function
func LoginOrganizerController(c echo.Context) error {
	login := models.LoginRequestBody{}
	// Bind all data from JSON
	c.Bind(&login)
	organizer, err := database.LoginOrganizer(login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	if organizer == nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid email or password"))
	}
	token, err := middlewares.CreateToken(int(organizer.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("can not generate token"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccessLogin("login success", organizer.ID, token, organizer.WoName, "organizer"))
}

// Get Profile Organizer Function
func GetProfileOrganizerController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	respon, err := database.FindProfilOrganizer(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get organizer", respon))
}

// Update/Edit Profile Organizer Function
func UpdateOrganizerController(c echo.Context) error {
	organizer := models.Organizer{}
	c.Bind(&organizer)
	organizer_id := middlewares.ExtractTokenUserId(c)
	_, err := database.EditOrganizer(organizer, organizer_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success edit data"))
}

// Update/Edit Profil Photo Organizer Function
func UpdatePhotoOrganizerController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	dataWo, _ := database.FindOrganizerById(organizer_id)
	// Process Upload Photo to Google Cloud
	bucket := "alta_wedding"
	var err error
	ctx := appengine.NewContext(c.Request())
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	f, uploaded_file, err := c.Request().FormFile("logo")
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
	organizerName := strings.ReplaceAll(dataWo.WoName, " ", "+")
	uploaded_file.Filename = fmt.Sprintf("%v-%v.%v", organizerName, formatted, extension)
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
	urlLogo := fmt.Sprintf("%v", u)
	_, e := database.EditPhotoOrganizer(urlLogo, organizer_id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success upload photo"))
}

// Testing Get User
func GetProfileOrganizerControllerTest() echo.HandlerFunc {
	return GetProfileOrganizerController
}
