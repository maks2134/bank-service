package utils

import (
	"bank_service/pkg/db"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

type ExportManager struct{}

func NewExportManager() *ExportManager {
	return &ExportManager{}
}

func (em *ExportManager) ExportTableToExcel(tableName string, filename string) error {
	rows, err := db.DB.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Write headers
	for i, col := range columns {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, col)
	}

	// Write data
	rowNum := 2
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, val := range values {
			cell := fmt.Sprintf("%c%d", 'A'+i, rowNum)
			if val != nil {
				f.SetCellValue(sheetName, cell, val)
			}
		}
		rowNum++
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

func (em *ExportManager) ExportQueryResultsToExcel(query string, filename string) error {
	rows, err := db.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Write headers
	for i, col := range columns {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, col)
	}

	// Write data
	rowNum := 2
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, val := range values {
			cell := fmt.Sprintf("%c%d", 'A'+i, rowNum)
			if val != nil {
				f.SetCellValue(sheetName, cell, val)
			}
		}
		rowNum++
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

func (em *ExportManager) ExportTableToCSV(tableName string, filename string) error {
	rows, err := db.DB.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write headers
	for i, col := range columns {
		if i > 0 {
			file.WriteString(",")
		}
		file.WriteString(fmt.Sprintf(`"%s"`, col))
	}
	file.WriteString("\n")

	// Write data
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, val := range values {
			if i > 0 {
				file.WriteString(",")
			}
			if val != nil {
				file.WriteString(fmt.Sprintf(`"%v"`, val))
			} else {
				file.WriteString(`""`)
			}
		}
		file.WriteString("\n")
	}

	return nil
}
