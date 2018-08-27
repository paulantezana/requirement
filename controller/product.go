package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
)

func GetProducts(c echo.Context) error {
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
	products := make([]models.Product, 0)

	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Order("id desc").
		Offset(offset).Limit(request.Limit).Find(&products).
		Offset(-1).Limit(-1).Count(&total).
		Error; err != nil {
		return err
	}

	// Type response
	// 0 = all data
	// 1 = minimal data
	if request.Type == 1 {
		customProducts := make([]models.Product, 0)
		for _, product := range products {
			customProducts = append(customProducts, models.Product{
				ID:   product.ID,
				Name: product.Name,
			})
		}
		return c.JSON(http.StatusCreated, utilities.Response{
			Success:     true,
			Data:        customProducts,
			Total:       total,
			CurrentPage: request.CurrentPage,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success:     true,
		Data:        products,
		Total:       total,
		CurrentPage: request.CurrentPage,
	})
}

func GetProductSearch(c echo.Context) error {
	// Get data request
	request := utilities.Request{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	products := make([]models.Product, 0)
	if err := db.Where("lower(name) LIKE lower(?)", "%"+request.Search+"%").
		Limit(5).Find(&products).Error; err != nil {
		return err
	}

	customProducts := make([]models.Product, 0)
	for _, product := range products {
		customProducts = append(customProducts, models.Product{
			ID:   product.ID,
			Name: product.Name,
			UnitMeasure: product.UnitMeasure,
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    customProducts,
	})
}

func GetProductByID(c echo.Context) error {
	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return err
	}

	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Execute instructions
	if err := db.First(&product, product.ID).Error; err != nil {
		return err
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product,
	})
}

func CreateProduct(c echo.Context) error {
	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Insert product in database
	if err := db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product.ID,
		Message: fmt.Sprintf("El producto %s se registro exitosamente", product.Name),
	})
}

func UpdateProduct(c echo.Context) error {
	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Update product in database
	rows := db.Model(&product).Update(product).RowsAffected
	if !product.State {
		rows = db.Model(product).UpdateColumn("state", false).RowsAffected
	}
	if rows == 0 {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se pudo actualizar el registro con el id = %d", product.ID),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product.ID,
	})
}

func DeleteProduct(c echo.Context) error {
	// Get data request
	product := models.Product{}
	if err := c.Bind(&product); err != nil {
		return err
	}

	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Validation product exist
	if db.First(&product).RecordNotFound() {
		return c.JSON(http.StatusCreated, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("No se encontró el registro con id %d", product.ID),
		})
	}

	// Delete product in database
	if err := db.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusOK, utilities.Response{
			Success: false,
			Message: fmt.Sprintf("%s", err),
		})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    product.ID,
	})
}
