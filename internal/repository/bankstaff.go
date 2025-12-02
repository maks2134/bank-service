package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
	"database/sql"
)

type BankStaffRepository struct{}

func NewBankStaffRepository() *BankStaffRepository {
	return &BankStaffRepository{}
}

func (r *BankStaffRepository) GetAll() ([]models.BankStaff, error) {
	rows, err := db.DB.Query("SELECT id, fullname, passport, position, hiredate, accesslevel, qualification, branchid FROM bankstaff ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staff []models.BankStaff
	for rows.Next() {
		var s models.BankStaff
		var qual sql.NullString
		err := rows.Scan(&s.ID, &s.FullName, &s.Passport, &s.Position, &s.HireDate, &s.AccessLevel, &qual, &s.BranchID)
		if err != nil {
			return nil, err
		}
		if qual.Valid {
			s.Qualification = &qual.String
		}
		staff = append(staff, s)
	}
	return staff, nil
}

func (r *BankStaffRepository) GetByID(id int) (*models.BankStaff, error) {
	var s models.BankStaff
	var qual sql.NullString
	err := db.DB.QueryRow("SELECT id, fullname, passport, position, hiredate, accesslevel, qualification, branchid FROM bankstaff WHERE id = $1", id).
		Scan(&s.ID, &s.FullName, &s.Passport, &s.Position, &s.HireDate, &s.AccessLevel, &qual, &s.BranchID)
	if err != nil {
		return nil, err
	}
	if qual.Valid {
		s.Qualification = &qual.String
	}
	return &s, nil
}

func (r *BankStaffRepository) Create(s *models.BankStaff) error {
	err := db.DB.QueryRow(
		"INSERT INTO bankstaff (fullname, passport, position, hiredate, accesslevel, qualification, branchid) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		s.FullName, s.Passport, s.Position, s.HireDate, s.AccessLevel, s.Qualification, s.BranchID,
	).Scan(&s.ID)
	return err
}

func (r *BankStaffRepository) Update(s *models.BankStaff) error {
	_, err := db.DB.Exec(
		"UPDATE bankstaff SET fullname=$1, passport=$2, position=$3, hiredate=$4, accesslevel=$5, qualification=$6, branchid=$7 WHERE id=$8",
		s.FullName, s.Passport, s.Position, s.HireDate, s.AccessLevel, s.Qualification, s.BranchID, s.ID,
	)
	return err
}

func (r *BankStaffRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM bankstaff WHERE id = $1", id)
	return err
}
