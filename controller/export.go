package controller

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/paulantezana/requirement/config"
	"github.com/paulantezana/requirement/models"
)

func ExportRequirementAll(c echo.Context) error {
	// get connection
	db := config.GetConnection()
	defer db.Close()

	// Query get config app
	con := models.Setting{}
	db.First(&con)

	// Create new BOOK EXCEL
	xlsx := excelize.NewFile()

	err := xlsx.AddPicture("Sheet1", "B2", "./static/logo.png", `{"x_scale": 0.5, "y_scale": 0.5}`)
	if err != nil {
		return err
	}
	xlsx.SetCellValue("Sheet1", "A5", con.CompanyName)
	xlsx.SetCellValue("Sheet1", "A6", con.City)
	xlsx.SetCellValue("Sheet1", "A8", "Requerimiento")

	//SET HEADER TABLE
	xlsx.SetCellValue("Sheet1", "A10", "Requerimiento")
	xlsx.SetCellValue("Sheet1", "B10", "Lugar")
	xlsx.SetCellValue("Sheet1", "C10", "Destino")
	xlsx.SetCellValue("Sheet1", "D10", "Fecha Emision")
	xlsx.SetCellValue("Sheet1", "E10", "Estado")

	// Get all requirements
	requirements := make([]models.Requirement, 0)
	db.Find(&requirements)

	currentRow := 11
	for k, rq := range requirements {
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow+k), rq.Name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow+k), rq.Place)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow+k), rq.Destination)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("D%d", currentRow+k), rq.EmissionDate)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("E%d", currentRow+k), rq.State)
	}

	fileAddress := "templates/requerimeinto.xlsx"

	err = xlsx.SaveAs("./" + fileAddress)
	if err != nil {
		return err
	}

	return c.File(fileAddress)
}
