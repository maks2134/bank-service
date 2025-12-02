package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
	"database/sql"
)

type CreditRepository struct{}

func NewCreditRepository() *CreditRepository {
	return &CreditRepository{}
}

func (r *CreditRepository) GetAll() ([]models.Credit, error) {
	rows, err := db.DB.Query("SELECT id, amount, interestrate, purpose, issuedate, status, repaymentdate, currency, risk_level FROM credit ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []models.Credit
	for rows.Next() {
		var c models.Credit
		var repaymentDate sql.NullString
		var riskLevel sql.NullString
		err := rows.Scan(&c.ID, &c.Amount, &c.InterestRate, &c.Purpose, &c.IssueDate, &c.Status, &repaymentDate, &c.Currency, &riskLevel)
		if err != nil {
			return nil, err
		}
		if repaymentDate.Valid {
			c.RepaymentDate = &repaymentDate.String
		}
		if riskLevel.Valid {
			c.RiskLevel = &riskLevel.String
		}
		credits = append(credits, c)
	}
	return credits, nil
}

func (r *CreditRepository) GetByID(id int) (*models.Credit, error) {
	var c models.Credit
	var repaymentDate sql.NullString
	var riskLevel sql.NullString
	err := db.DB.QueryRow("SELECT id, amount, interestrate, purpose, issuedate, status, repaymentdate, currency, risk_level FROM credit WHERE id = $1", id).
		Scan(&c.ID, &c.Amount, &c.InterestRate, &c.Purpose, &c.IssueDate, &c.Status, &repaymentDate, &c.Currency, &riskLevel)
	if err != nil {
		return nil, err
	}
	if repaymentDate.Valid {
		c.RepaymentDate = &repaymentDate.String
	}
	if riskLevel.Valid {
		c.RiskLevel = &riskLevel.String
	}
	return &c, nil
}

func (r *CreditRepository) Create(c *models.Credit) error {
	err := db.DB.QueryRow(
		"INSERT INTO credit (amount, interestrate, purpose, issuedate, status, repaymentdate, currency, risk_level) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		c.Amount, c.InterestRate, c.Purpose, c.IssueDate, c.Status, c.RepaymentDate, c.Currency, c.RiskLevel,
	).Scan(&c.ID)
	return err
}

func (r *CreditRepository) Update(c *models.Credit) error {
	_, err := db.DB.Exec(
		"UPDATE credit SET amount=$1, interestrate=$2, purpose=$3, issuedate=$4, status=$5, repaymentdate=$6, currency=$7, risk_level=$8 WHERE id=$9",
		c.Amount, c.InterestRate, c.Purpose, c.IssueDate, c.Status, c.RepaymentDate, c.Currency, c.RiskLevel, c.ID,
	)
	return err
}

func (r *CreditRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM credit WHERE id = $1", id)
	return err
}
