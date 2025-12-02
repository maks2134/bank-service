package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
)

type CreditBankStaffRepository struct{}

func NewCreditBankStaffRepository() *CreditBankStaffRepository {
	return &CreditBankStaffRepository{}
}

func (r *CreditBankStaffRepository) GetAll() ([]models.CreditBankStaff, error) {
	rows, err := db.DB.Query("SELECT creditid, staffid FROM creditbankstaff ORDER BY creditid, staffid")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []models.CreditBankStaff
	for rows.Next() {
		var rel models.CreditBankStaff
		err := rows.Scan(&rel.CreditID, &rel.StaffID)
		if err != nil {
			return nil, err
		}
		relations = append(relations, rel)
	}
	return relations, nil
}

func (r *CreditBankStaffRepository) Create(rel *models.CreditBankStaff) error {
	_, err := db.DB.Exec(
		"INSERT INTO creditbankstaff (creditid, staffid) VALUES ($1, $2) ON CONFLICT (creditid, staffid) DO NOTHING",
		rel.CreditID, rel.StaffID,
	)
	return err
}

func (r *CreditBankStaffRepository) Delete(creditID, staffID int) error {
	_, err := db.DB.Exec("DELETE FROM creditbankstaff WHERE creditid = $1 AND staffid = $2", creditID, staffID)
	return err
}
