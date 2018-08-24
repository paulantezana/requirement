package controller

import (
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
)

type GlobalSettings struct {
	Messages []string       `json:"messages"`
	Success  bool           `json:"success"`
	Setting  models.Setting `json:"setting"`
	User     models.User    `json:"user"`
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
		User:     user,
		Setting:  con,
		Success:  true,
		Messages: []string{"ok"},
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
	})
}
