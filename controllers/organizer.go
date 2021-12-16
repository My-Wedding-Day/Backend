package controllers

import (
	"alta-wedding/lib/database"
	"alta-wedding/lib/responses"
	"alta-wedding/middlewares"
	"alta-wedding/models"
	"alta-wedding/util"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
	// Check data cannot be empty
	if organizer.Email == "" || organizer.City == "" {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("input data cannot be empty"))
	}
	// Check Name cannot less than 5 characters
	if len(organizer.WoName) < 5 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("business name cannot less than 5 characters"))
	}
	// Check Organizer Email is Exist
	emailCheck, _ := database.CheckDatabase("email", organizer.Email)
	if emailCheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email was used, try another one"))
	}
	// Check Organizer Business name is Exist
	nameCheck, _ := database.CheckDatabase("wo_name", organizer.WoName)
	if nameCheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("business name was used, try another one"))
	}
	// REGEX
	var pattern string
	var matched bool
	// Check Format Name
	pattern = `^\w(\w+ ?)*$`
	regex, _ := regexp.Compile(pattern)
	matched = regex.Match([]byte(organizer.WoName))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("invalid format name"))
	}
	// Check Format Email
	pattern = `^\w+@\w+\.\w+$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Email))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("email must contain email format"))
	}
	// Check Format Phone Number
	pattern = `^[0-9]*$`
	matched, _ = regexp.Match(pattern, []byte(organizer.PhoneNumber))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone must be number"))
	}
	// Check Format Address
	pattern = `^[a-zA-Z]([a-zA-Z.0-9,]+ ?)*$`
	matched, _ = regexp.Match(pattern, []byte(organizer.Address))
	if !matched {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Address must be valid"))
	}
	// Check Address
	_, _, Err := util.GetGeocodeLocations(organizer.Address)
	if Err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("Address "+Err.Error()))
	}
	// Check Length of Character of PhoneNumber and Password
	if len(organizer.PhoneNumber) < 9 || len(organizer.Password) < 8 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("password or phone number cannot less than 8 characters"))
	}
	// Check Phone number existing
	phonecheck, _ := database.CheckDatabase("phone_number", organizer.PhoneNumber)
	if phonecheck > 0 {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("phone number was used, try another one"))
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
	token, _ := middlewares.CreateToken(int(organizer.ID))
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

// Get Profile Organizer by ID
func GetProileOrganizerbyIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
	}
	respon, err := database.FindProfilOrganizer(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get organizer", respon))
}

// Get my Package for Organizer
func GetMyPackageController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	mypackages, err := database.GetPackagesByToken(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get my packages", mypackages))
}

// Get My Reservation List From Users Order
func GetMyReservationListController(c echo.Context) error {
	organizer_id := middlewares.ExtractTokenUserId(c)
	mylistorder, err := database.GetListReservations(organizer_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusOK, responses.StatusSuccessData("success get my list order", mylistorder))
}

// Fitur Accept/Decline Reservation
// func AcceptDeclineController(c echo.Context) error {
// 	reservation_id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("false param"))
// 	}
// 	request := models.AcceptBody{}
// 	c.Bind(&request)
// 	if request.Status_Order != "accept" || request.Status_Order != "decline" {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
// 	}
// 	_, err := database.AcceptDecline(reservation_id, request.Status_Order)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, responses.StatusFailed("bad request"))
// 	}
// 	return c.JSON(http.StatusCreated, responses.StatusSuccess("success edit data"))
// }

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

// Insert Document Organizer Function
func UpdateDocumentsOrganizerController(c echo.Context) error {
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
	f, uploaded_file, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailedDataPhoto(err.Error()))
	}
	defer f.Close()
	fileExtensions := map[string]bool{"jpg": true, "jpeg": true, "png": true, "pdf": true}
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
	urlDocument := fmt.Sprintf("%v", u)
	_, e := database.EditDocumentOrganizer(urlDocument, organizer_id)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, responses.StatusFailed("internal server error"))
	}
	return c.JSON(http.StatusCreated, responses.StatusSuccess("success upload document"))
}

// Testing Get User
func GetProfileOrganizerControllerTest() echo.HandlerFunc {
	return GetProfileOrganizerController
}
