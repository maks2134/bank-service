package utils

import (
	"bank_service/pkg/db"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type BackupManager struct{}

func NewBackupManager() *BackupManager {
	return &BackupManager{}
}

type BackupData struct {
	Timestamp time.Time              `json:"timestamp"`
	Tables    map[string]interface{} `json:"tables"`
}

func (bm *BackupManager) BackupTable(tableName string) (string, error) {
	data, columns, err := NewTableManager().GetTableData(tableName)
	if err != nil {
		return "", err
	}

	backup := BackupData{
		Timestamp: time.Now(),
		Tables: map[string]interface{}{
			tableName: map[string]interface{}{
				"columns": columns,
				"data":    data,
			},
		},
	}

	filename := fmt.Sprintf("backup_%s_%s.json", tableName, time.Now().Format("20060102_150405"))
	dataBytes, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(filename, dataBytes, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (bm *BackupManager) BackupDatabase() (string, error) {
	tm := NewTableManager()
	tables, err := tm.GetTables()
	if err != nil {
		return "", err
	}

	backup := BackupData{
		Timestamp: time.Now(),
		Tables:    make(map[string]interface{}),
	}

	for _, tableName := range tables {
		data, columns, err := tm.GetTableData(tableName)
		if err != nil {
			return "", err
		}
		backup.Tables[tableName] = map[string]interface{}{
			"columns": columns,
			"data":    data,
		}
	}

	filename := fmt.Sprintf("backup_db_%s.json", time.Now().Format("20060102_150405"))
	dataBytes, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(filename, dataBytes, 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (bm *BackupManager) RestoreTable(filename string, tableName string) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var backup BackupData
	if err := json.Unmarshal(dataBytes, &backup); err != nil {
		return err
	}

	tableData, ok := backup.Tables[tableName].(map[string]interface{})
	if !ok {
		return fmt.Errorf("table %s not found in backup", tableName)
	}

	columns, ok := tableData["columns"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid backup format")
	}

	// Drop existing table if exists
	_ = NewTableManager().DropTable(tableName)

	// Recreate table structure
	colDefs := make([]map[string]string, len(columns))
	for i, col := range columns {
		colDefs[i] = map[string]string{
			"name": col.(string),
			"type": "text", // Default type for restore
		}
	}

	if err := NewTableManager().CreateTable(tableName, colDefs); err != nil {
		return err
	}

	// Restore data
	data, ok := tableData["data"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid backup format")
	}

	for _, row := range data {
		values := row.([]interface{})
		placeholders := ""
		for i := range values {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += fmt.Sprintf("$%d", i+1)
		}

		colNames := ""
		for i, col := range columns {
			if i > 0 {
				colNames += ", "
			}
			colNames += col.(string)
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, colNames, placeholders)
		_, err := db.DB.Exec(query, values...)
		if err != nil {
			return err
		}
	}

	return nil
}
