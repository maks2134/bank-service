package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
	"database/sql"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) GetAll() ([]models.Transaction, error) {
	rows, err := db.DB.Query("SELECT id, amount, operationdate, operationtype, purpose, operationstatus, currency, branchid, accountid FROM transaction ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		var purpose sql.NullString
		err := rows.Scan(&t.ID, &t.Amount, &t.OperationDate, &t.OperationType, &purpose, &t.OperationStatus, &t.Currency, &t.BranchID, &t.AccountID)
		if err != nil {
			return nil, err
		}
		if purpose.Valid {
			t.Purpose = &purpose.String
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByID(id int) (*models.Transaction, error) {
	var t models.Transaction
	var purpose sql.NullString
	err := db.DB.QueryRow("SELECT id, amount, operationdate, operationtype, purpose, operationstatus, currency, branchid, accountid FROM transaction WHERE id = $1", id).
		Scan(&t.ID, &t.Amount, &t.OperationDate, &t.OperationType, &purpose, &t.OperationStatus, &t.Currency, &t.BranchID, &t.AccountID)
	if err != nil {
		return nil, err
	}
	if purpose.Valid {
		t.Purpose = &purpose.String
	}
	return &t, nil
}

func (r *TransactionRepository) Create(t *models.Transaction) error {
	err := db.DB.QueryRow(
		"INSERT INTO transaction (amount, operationdate, operationtype, purpose, operationstatus, currency, branchid, accountid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		t.Amount, t.OperationDate, t.OperationType, t.Purpose, t.OperationStatus, t.Currency, t.BranchID, t.AccountID,
	).Scan(&t.ID)
	return err
}

func (r *TransactionRepository) Update(t *models.Transaction) error {
	_, err := db.DB.Exec(
		"UPDATE transaction SET amount=$1, operationdate=$2, operationtype=$3, purpose=$4, operationstatus=$5, currency=$6, branchid=$7, accountid=$8 WHERE id=$9",
		t.Amount, t.OperationDate, t.OperationType, t.Purpose, t.OperationStatus, t.Currency, t.BranchID, t.AccountID, t.ID,
	)
	return err
}

func (r *TransactionRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM transaction WHERE id = $1", id)
	return err
}
