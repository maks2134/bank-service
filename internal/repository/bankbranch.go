package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
)

type BankBranchRepository struct{}

func NewBankBranchRepository() *BankBranchRepository {
	return &BankBranchRepository{}
}

func (r *BankBranchRepository) GetAll() ([]models.BankBranch, error) {
	rows, err := db.DB.Query("SELECT id, servicezone, full_address, contactphone, workingdays, workinghours, branchtype FROM bankbranch ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []models.BankBranch
	for rows.Next() {
		var branch models.BankBranch
		err := rows.Scan(&branch.ID, &branch.ServiceZone, &branch.FullAddress, &branch.ContactPhone, &branch.WorkingDays, &branch.WorkingHours, &branch.BranchType)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

func (r *BankBranchRepository) GetByID(id int) (*models.BankBranch, error) {
	var branch models.BankBranch
	err := db.DB.QueryRow("SELECT id, servicezone, full_address, contactphone, workingdays, workinghours, branchtype FROM bankbranch WHERE id = $1", id).
		Scan(&branch.ID, &branch.ServiceZone, &branch.FullAddress, &branch.ContactPhone, &branch.WorkingDays, &branch.WorkingHours, &branch.BranchType)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *BankBranchRepository) Create(branch *models.BankBranch) error {
	err := db.DB.QueryRow(
		"INSERT INTO bankbranch (servicezone, full_address, contactphone, workingdays, workinghours, branchtype) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		branch.ServiceZone, branch.FullAddress, branch.ContactPhone, branch.WorkingDays, branch.WorkingHours, branch.BranchType,
	).Scan(&branch.ID)
	return err
}

func (r *BankBranchRepository) Update(branch *models.BankBranch) error {
	_, err := db.DB.Exec(
		"UPDATE bankbranch SET servicezone=$1, full_address=$2, contactphone=$3, workingdays=$4, workinghours=$5, branchtype=$6 WHERE id=$7",
		branch.ServiceZone, branch.FullAddress, branch.ContactPhone, branch.WorkingDays, branch.WorkingHours, branch.BranchType, branch.ID,
	)
	return err
}

func (r *BankBranchRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM bankbranch WHERE id = $1", id)
	return err
}
