package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
    "github.com/paulantezana/requirement/api"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	migration()
	//initial()
	//TestBorrame()

	// COR
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// Static Files =========================================================================
	static := e.Group("/static")
	static.Static("", "static")

	// API
	api.PublicApi(e)
	api.ProtectedApi(e)

	// Custom port
	port := os.Getenv("PORT")
	if port == "" {
		port = config.GetConfig().Server.Port
	}

	// Starting server echo
	e.Logger.Fatal(e.Start(":" + port))
}

func TestBorrame() {
	fmt.Println("holaaaaaaaaaaaaa")
}

func migration() {
	db := config.GetConnection()
	defer db.Close()

	db.Debug().AutoMigrate(
		&models.User{},
		&models.Quotation{},
		&models.QuotationDetail{},
		&models.Product{},
		&models.Provider{},
		&models.Requirement{},
		&models.Require{},
		&models.Setting{},
	)
	db.Model(&models.Requirement{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Require{}).AddForeignKey("requirement_id", "requirements(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Require{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.QuotationDetail{}).AddForeignKey("require_id", "requires(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.QuotationDetail{}).AddForeignKey("quotation_id", "quotations(id)", "RESTRICT", "RESTRICT")

	db.Model(&models.Quotation{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Quotation{}).AddForeignKey("provider_id", "providers(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Quotation{}).AddForeignKey("requirement_id", "requirements(id)", "RESTRICT", "RESTRICT")
}
func initial() {
	db := config.GetConnection()
	defer db.Close()

	conf := models.Setting{
		Company:        "ABC Company",
		City:           "Sicuani - Cusco - Perú",
		Identification: "20563258412",
		Email:          "empresa@empresa.com",
		Item:           10,
		Quotations:     3,
	}

	if err := db.Create(&conf).Error; err != nil {
		log.Fatal(err)
	}
}
