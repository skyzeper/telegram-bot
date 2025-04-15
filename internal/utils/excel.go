package utils

import (
	"github.com/xuri/excelize/v2"
)

func CreateExcelFile(filename string) (*excelize.File, error) {
	f := excelize.NewFile()
	return f, f.SaveAs(filename)
}
