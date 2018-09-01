package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"io"
	"net/http"
	"os"
	"path"
)

type GlobalSettings struct {
	Message string         `json:"message"`
	Success bool           `json:"success"`
	Setting models.Setting `json:"setting"`
	User    models.User    `json:"user"`
}

func GetGlobalSettings(c echo.Context) error {
	// Get data request
	con := models.Setting{}
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&user, user.ID).Error; err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	user.Password = ""
	user.Key = ""

	db.First(&con) // Find settings

	// Set object response
	return c.JSON(http.StatusOK, GlobalSettings{
		User:    user,
		Setting: con,
		Success: true,
		Message: "OK",
	})
}

func GetSetting(c echo.Context) error {
	// Declare variables
	con := models.Setting{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database
	db.First(&con)

	// Response config
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    con,
	})
}

func UpdateSetting(c echo.Context) error {
	// Get data request
	con := models.Setting{}
	if err := c.Bind(&con); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation first data
	var exist uint
	db.Model(&models.Setting{}).Count(&exist)

	// Insert config in database
	if exist == 0 {
		if err := db.Create(&con).Error; err != nil {
			return err
		}
	}

	// Update con in database
	if err := db.Model(&con).Update(con).Error; err != nil {
		return err
	}

	// Response config
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    con.ID,
		Message: "OK",
	})
}

func UploadLogoSetting(c echo.Context) error {
	// Read form fields
	idSetting := c.FormValue("id")
	setting := models.Setting{}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	if db.First(&setting, "id = ?", idSetting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}

	// Source
	file, err := c.FormFile("logo")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	logoSRC := fmt.Sprintf("static/logo%s", path.Ext(file.Filename))
	dst, err := os.Create(logoSRC)
	if err != nil {
		return err
	}
	defer dst.Close()
	setting.Logo = logoSRC

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update database user
	if err := db.Model(&setting).Update(setting).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    setting.ID,
		Message: "OK",
	})
}

func DownloadLogoSetting(c echo.Context) error {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation user exist
	setting := models.Setting{}
	if db.First(&setting).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", setting.ID),
		})
	}
	return c.File(setting.Logo)
}
