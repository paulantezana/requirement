package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
)

func GetProviders(c echo.Context) error {
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
	providers := make([]models.Provider, 0)

	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Or("ruc LIKE ?", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&providers).
		Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customProvider := make([]models.Provider, 0)
		for _, provider := range providers {
			customProvider = append(customProvider, models.Provider{
				ID:   provider.ID,
				Name: provider.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.Response{
			Success:     true,
			Data:        customProvider,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        providers,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetProviderSearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	providers := make([]models.Provider, 0)
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Or("ruc LIKE ?", "%"+request.Search+"%").
		Limit(5).Find(&providers).Error; err != nil {
		return err
	}

	customProviders := make([]models.Product, 0)
	for _, product := range providers {
		customProviders = append(customProviders, models.Product{
			ID:   product.ID,
			Name: product.Name,
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customProviders,
	})
}

func GetProviderByID(c echo.Context) error {
	// Get data request
	provider := models.Provider{}
	if err := c.Bind(&provider); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&provider, provider.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    provider,
	})
}

func CreateProvider(c echo.Context) error {
	// Get data request
	provider := models.Provider{}
	if err := c.Bind(&provider); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert provider in database
	if err := db.Create(&provider).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    provider.ID,
		Message: fmt.Sprintf("El proveedor %s se registro exitosamente", provider.Name),
	})
}

func UpdateProvider(c echo.Context) error {
	// Get data request
	provider := models.Provider{}
	if err := c.Bind(&provider); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update provider in database
	rows := db.Model(&provider).Update(provider).RowsAffected
	if !provider.State {
		rows = db.Model(provider).UpdateColumn("state", false).RowsAffected
	}
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", provider.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    provider.ID,
	})
}

func DeleteProvider(c echo.Context) error {
	// Get data request
	provider := models.Provider{}
	if err := c.Bind(&provider); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation provider exist
	if db.First(&provider).RecordNotFound() {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", provider.ID),
		})
	}

	// Delete provider in database
	if err := db.Delete(&provider).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    provider.ID,
	})
}
