package controller

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
		return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
			Success:     true,
			Data:        customProvider,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.ResponsePaginate{
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
			Message: fmt.Sprintf("No se encontró el registro con id %d", provider.ID),
		})
	}

	// Delete provider in database
	if err := db.Delete(&provider).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    provider.ID,
	})
}

func ValidateRucProvider(c echo.Context) error {
	// Get data request
	provider := models.Provider{}
	if err := c.Bind(&provider); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validations
	if err := db.Where("ruc = ?", provider.RUC).First(&provider).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: true,
			Message: "OK",
		})
	}

	// Return response
	return c.JSON(http.StatusOK, utilities.Response{
		Success: false,
		Message: fmt.Sprintf("El número de RUC ya esta registrado"),
	})
}

func GetTempUploadProvider(c echo.Context) error {
	return c.File("templates/uploadProviderTemplate.xlsx")
}

func SetTempUploadProvider(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	auxDir := "temp/provider" + filepath.Ext(file.Filename)
	dst, err := os.Create(auxDir)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// ---------------------
	// Read File whit Excel
	// ---------------------
	xlsx, err := excelize.OpenFile(auxDir)
	if err != nil {
		return err
	}

	// Prepare
	providers := make([]models.Provider, 0)
	ignoreCols := 1

	// Get all the rows in the proveedores(Sheet).
	rows := xlsx.GetRows("proveedores")
	for k, row := range rows {
		if k >= ignoreCols {
			providers = append(providers, models.Provider{
				RUC:     strings.TrimSpace(row[0]),
				Name:    strings.TrimSpace(row[1]),
				Manager: strings.TrimSpace(row[5]),
				Email:   strings.TrimSpace(row[5]),
				Phone:   strings.TrimSpace(row[5]),
				Address: strings.TrimSpace(row[5]),
				State:   true,
			})
		}
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert providers in database
	tr := db.Begin()
	for _, provider := range providers {
		if err := tr.Create(&provider).Error; err != nil {
			tr.Rollback()
			return c.JSON(http.StatusOK, utilities.Response{
				Success: false,
				Message: fmt.Sprintf("Ocurrió un error al insertar el proveedores %s con "+
					"RUC: %s es posible que este proveedor ya este en la base de datos o los datos son incorrectos, "+
					"Error: %s, no se realizo ninguna cambio en la base de datos", provider.Name, provider.RUC, err),
			})
		}
	}
	tr.Commit()

	// Response success
	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Message: fmt.Sprintf("Se guardo %d registros den la base de datos", len(providers)),
	})
}
