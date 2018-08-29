package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
	"time"
)

func GetRequirements(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Pagination calculate
	if request.CurrentPage == 0 {
		request.CurrentPage = 1
	}
	offset := request.Limit*request.CurrentPage - request.Limit

	// Execute instructions
	var total uint
	requirements := make([]models.Requirement, 0)

	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(destination) LIKE lower(?)", "%"+request.Search+"%").
		Or("lower(claimant) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&requirements).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customRequirements := make([]models.Requirement, 0)
		for _, product := range requirements {
			customRequirements = append(customRequirements, models.Requirement{
				ID:   product.ID,
				Name: product.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customRequirements,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
		Success:     true,
		Data:        requirements,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetRequirementByID(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&requirement, requirement.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement,
	})
}

func CreateRequirement(c echo.Context) error {
	// Get user token authenticate
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utilities.Claim)
	currentUser := claims.User

	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}
	requirement.UserID = currentUser.ID

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation
	if len(requirement.Requires) == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: "Agregue al menos un producto para crear el requerimiento",
		})
	}

	// Default values
	requirement.EmissionDate = time.Now()
	requirement.State = "0"

	// Insert product in database
	if err := db.Create(&requirement).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement.ID,
		Message: fmt.Sprintf("El requerimiento %s se registro exitosamente", requirement.Name),
	})
}

func UpdateRequirement(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update product in database
	rows := db.Debug().Model(&requirement).Update(requirement).RowsAffected

	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", requirement.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement.ID,
	})
}

func SetRejectedRequirement(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update product in database
	rows := db.Model(&requirement).UpdateColumn("state", "2").RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", requirement.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement.ID,
	})
}

func SetClosedRequirement(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update product in database
	rows := db.Model(&requirement).UpdateColumn("state", "4").RowsAffected
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", requirement.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement.ID,
	})
}

func DeleteRequirement(c echo.Context) error {
	// Get data request
	requirement := models.Requirement{}
	if err := c.Bind(&requirement); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation product exist
	if db.First(&requirement).RecordNotFound() {
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", requirement.ID),
		})
	}

	// Delete product in database
	if err := db.Delete(&requirement).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    requirement.ID,
	})
}
