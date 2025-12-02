package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
	"database/sql"
)

type BankAccountRepository struct{}

func NewBankAccountRepository() *BankAccountRepository {
	return &BankAccountRepository{}
}

func (r *BankAccountRepository) GetAll() ([]models.BankAccount, error) {
	rows, err := db.DB.Query("SELECT id, accounttype, accountnumber, balance, currency, opendate, accountstatus, clientid FROM bankaccount ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.BankAccount
	for rows.Next() {
		var acc models.BankAccount
		var clientID sql.NullInt64
		err := rows.Scan(&acc.ID, &acc.AccountType, &acc.AccountNumber, &acc.Balance, &acc.Currency, &acc.OpenDate, &acc.AccountStatus, &clientID)
		if err != nil {
			return nil, err
		}
		if clientID.Valid {
			id := int(clientID.Int64)
			acc.ClientID = &id
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *BankAccountRepository) GetByID(id int) (*models.BankAccount, error) {
	var acc models.BankAccount
	var clientID sql.NullInt64
	err := db.DB.QueryRow("SELECT id, accounttype, accountnumber, balance, currency, opendate, accountstatus, clientid FROM bankaccount WHERE id = $1", id).
		Scan(&acc.ID, &acc.AccountType, &acc.AccountNumber, &acc.Balance, &acc.Currency, &acc.OpenDate, &acc.AccountStatus, &clientID)
	if err != nil {
		return nil, err
	}
	if clientID.Valid {
		id := int(clientID.Int64)
		acc.ClientID = &id
	}
	return &acc, nil
}

func (r *BankAccountRepository) Create(acc *models.BankAccount) error {
	err := db.DB.QueryRow(
		"INSERT INTO bankaccount (accounttype, accountnumber, balance, currency, opendate, accountstatus, clientid) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		acc.AccountType, acc.AccountNumber, acc.Balance, acc.Currency, acc.OpenDate, acc.AccountStatus, acc.ClientID,
	).Scan(&acc.ID)
	return err
}

func (r *BankAccountRepository) Update(acc *models.BankAccount) error {
	_, err := db.DB.Exec(
		"UPDATE bankaccount SET accounttype=$1, accountnumber=$2, balance=$3, currency=$4, opendate=$5, accountstatus=$6, clientid=$7 WHERE id=$8",
		acc.AccountType, acc.AccountNumber, acc.Balance, acc.Currency, acc.OpenDate, acc.AccountStatus, acc.ClientID, acc.ID,
	)
	return err
}

func (r *BankAccountRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM bankaccount WHERE id = $1", id)
	return err
}
