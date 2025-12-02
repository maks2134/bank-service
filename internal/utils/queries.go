package utils

import (
	"bank_service/pkg/db"
	"encoding/json"
	"os"
	"time"
)

type QueryManager struct {
	queriesFile string
}

func NewQueryManager() *QueryManager {
	return &QueryManager{
		queriesFile: "saved_queries.json",
	}
}

type SavedQuery struct {
	Name        string    `json:"name"`
	Query       string    `json:"query"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type QueriesData struct {
	Queries []SavedQuery `json:"queries"`
}

func (qm *QueryManager) GetSavedQueries() ([]SavedQuery, error) {
	if _, err := os.Stat(qm.queriesFile); os.IsNotExist(err) {
		return []SavedQuery{}, nil
	}

	data, err := os.ReadFile(qm.queriesFile)
	if err != nil {
		return nil, err
	}

	var queriesData QueriesData
	if err := json.Unmarshal(data, &queriesData); err != nil {
		return nil, err
	}

	return queriesData.Queries, nil
}

func (qm *QueryManager) SaveQuery(query SavedQuery) error {
	queries, err := qm.GetSavedQueries()
	if err != nil {
		return err
	}

	query.CreatedAt = time.Now()
	queries = append(queries, query)

	queriesData := QueriesData{Queries: queries}
	data, err := json.MarshalIndent(queriesData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(qm.queriesFile, data, 0644)
}

func (qm *QueryManager) DeleteQuery(name string) error {
	queries, err := qm.GetSavedQueries()
	if err != nil {
		return err
	}

	var newQueries []SavedQuery
	for _, q := range queries {
		if q.Name != name {
			newQueries = append(newQueries, q)
		}
	}

	queriesData := QueriesData{Queries: newQueries}
	data, err := json.MarshalIndent(queriesData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(qm.queriesFile, data, 0644)
}

func (qm *QueryManager) ExecuteQuery(query string) ([][]interface{}, []string, error) {
	rows, err := db.DB.Query(query)
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

func (qm *QueryManager) GetPredefinedQueries() []SavedQuery {
	return []SavedQuery{
		{
			Name:        "Активные счета",
			Query:       "SELECT * FROM bankaccount WHERE accountstatus = 'ACTIVE'",
			Description: "Все активные банковские счета",
		},
		{
			Name:        "Клиенты с кредитами",
			Query:       "SELECT c.*, cr.amount, cr.purpose FROM client c JOIN credit cr ON c.creditid = cr.id",
			Description: "Клиенты и их кредиты",
		},
		{
			Name:        "Транзакции за месяц",
			Query:       "SELECT * FROM transaction WHERE operationdate >= CURRENT_DATE - INTERVAL '30 days'",
			Description: "Транзакции за последние 30 дней",
		},
		{
			Name:        "Сотрудники по отделениям",
			Query:       "SELECT bs.*, bb.servicezone FROM bankstaff bs JOIN bankbranch bb ON bs.branchid = bb.id",
			Description: "Сотрудники с информацией об отделениях",
		},
		{
			Name:        "Счета по валюте",
			Query:       "SELECT currency, COUNT(*) as count, SUM(balance) as total FROM bankaccount GROUP BY currency",
			Description: "Статистика по валютам",
		},
		{
			Name:        "Кредиты по статусу",
			Query:       "SELECT status, COUNT(*) as count, SUM(amount) as total FROM credit GROUP BY status",
			Description: "Статистика кредитов по статусам",
		},
	}
}
