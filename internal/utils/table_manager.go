package utils

import (
	"bank_service/pkg/db"
	"database/sql"
	"fmt"
)

type TableManager struct{}

func NewTableManager() *TableManager {
	return &TableManager{}
}

func (tm *TableManager) GetTables() ([]string, error) {
	rows, err := db.DB.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func (tm *TableManager) GetTableColumns(tableName string) ([]map[string]string, error) {
	rows, err := db.DB.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []map[string]string
	for rows.Next() {
		var colName, dataType, isNullable string
		var colDefault sql.NullString
		if err := rows.Scan(&colName, &dataType, &isNullable, &colDefault); err != nil {
			return nil, err
		}
		defaultValue := ""
		if colDefault.Valid {
			defaultValue = colDefault.String
		}
		columns = append(columns, map[string]string{
			"name":     colName,
			"type":     dataType,
			"nullable": isNullable,
			"default":  defaultValue,
		})
	}
	return columns, nil
}

func (tm *TableManager) CreateTable(tableName string, columns []map[string]string) error {
	query := fmt.Sprintf("CREATE TABLE %s (", tableName)
	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s %s", col["name"], col["type"])
		if col["nullable"] == "NO" {
			query += " NOT NULL"
		}
		if col["default"] != "" {
			query += fmt.Sprintf(" DEFAULT %s", col["default"])
		}
	}
	query += ")"

	_, err := db.DB.Exec(query)
	return err
}

func (tm *TableManager) DropTable(tableName string) error {
	_, err := db.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName))
	return err
}

func (tm *TableManager) AddColumn(tableName, columnName, dataType string, nullable bool) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, dataType)
	if !nullable {
		query += " NOT NULL"
	}
	_, err := db.DB.Exec(query)
	return err
}

func (tm *TableManager) DropColumn(tableName, columnName string) error {
	_, err := db.DB.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName))
	return err
}

func (tm *TableManager) GetTableData(tableName string) ([][]interface{}, []string, error) {
	rows, err := db.DB.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var data [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}
		data = append(data, values)
	}
	return data, columns, nil
}
