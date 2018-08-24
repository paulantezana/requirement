package controller

import (
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
	"github.com/paulantezana/requirement/utilities"
	"net/http"
)

type providerTop struct {
	ID     uint   `json:"-"`
	Winner bool   `json:"winner"`
	Name   string `json:"name"`
	Top    uint   `json:"top"`
}

func TopProviderWinner(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database top 5
	providerTops := make([]providerTop, 0)
	if err := db.Table("quotations").
		Select("providers.id, quotations.winner, providers.name, count(winner) as top").
		Joins("INNER JOIN providers on quotations.provider_id = providers.id").
		Group("providers.id,  providers.name, quotations.winner").
		Having("quotations.winner = true").
		Order("top desc").
		Limit(15).
		Scan(&providerTops).Error; err != nil {
		return err
	}

	// Total registers
	var total uint
	db.Model(models.Provider{}).Count(&total)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    providerTops,
		Total:   total,
	})
}

type userTop struct {
	ID        uint   `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Top       uint   `json:"top"`
}

func TopUsers(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database top 5
	userTops := make([]userTop, 0)
	if err := db.Table("quotations").
		Select("users.id, users.first_name, users.last_name, count(quotations.user_id) as top").
		Joins("INNER JOIN users on quotations.user_id = users.id").
		Group("users.id, users.first_name, users.last_name").
		Order("top desc").
		Limit(15).
		Scan(&userTops).Error; err != nil {
		return err
	}

	// Total registers
	var total uint
	db.Model(models.User{}).Count(&total)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    userTops,
		Total:   total,
	})
}

type productTop struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
	Top  uint   `json:"top"`
}

func TopProducts(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database top 5
	productTops := make([]productTop, 0)
	if err := db.Table("requires").
		Select("products.id, products.name, count(requires.product_id) as top").
		Joins("INNER JOIN products on requires.product_id = products.id").
		Group("products.id, products.name").
		Order("top desc").
		Limit(15).
		Scan(&productTops).Error; err != nil {
		return err
	}

	// Total registers
	var total uint
	db.Model(models.Product{}).Count(&total)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    productTops,
		Total:   total,
	})
}

type requirementStateTop struct {
	ID    uint   `json:"-"`
	State string `json:"state"`
	Top   uint   `json:"top"`
}

func TopRequirements(c echo.Context) error {
	// Get connection
	db := config.GetConnection()
	defer db.Close()

	// Query database top 5
	requirementStateTops := make([]requirementStateTop, 0)
	if err := db.Table("requirements").
		Select("state, count(*) as top").
		Group("state").
		Order("state").
		Scan(&requirementStateTops).Error; err != nil {
		return err
	}

	// Total registers
	var total uint
	db.Model(models.Requirement{}).Count(&total)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    requirementStateTops,
		Total:   total,
	})
}
