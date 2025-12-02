package repository

import (
	"bank_service/internal/models"
	"bank_service/pkg/db"
	"database/sql"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (r *ClientRepository) GetAll() ([]models.Client, error) {
	rows, err := db.DB.Query("SELECT id, phone, login, fullname, clienttype, password, email, passport, creditid FROM client ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		var creditID sql.NullInt64
		err := rows.Scan(&c.ID, &c.Phone, &c.Login, &c.FullName, &c.ClientType, &c.Password, &c.Email, &c.Passport, &creditID)
		if err != nil {
			return nil, err
		}
		if creditID.Valid {
			id := int(creditID.Int64)
			c.CreditID = &id
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (r *ClientRepository) GetByID(id int) (*models.Client, error) {
	var c models.Client
	var creditID sql.NullInt64
	err := db.DB.QueryRow("SELECT id, phone, login, fullname, clienttype, password, email, passport, creditid FROM client WHERE id = $1", id).
		Scan(&c.ID, &c.Phone, &c.Login, &c.FullName, &c.ClientType, &c.Password, &c.Email, &c.Passport, &creditID)
	if err != nil {
		return nil, err
	}
	if creditID.Valid {
		id := int(creditID.Int64)
		c.CreditID = &id
	}
	return &c, nil
}

func (r *ClientRepository) Create(c *models.Client) error {
	err := db.DB.QueryRow(
		"INSERT INTO client (phone, login, fullname, clienttype, password, email, passport, creditid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		c.Phone, c.Login, c.FullName, c.ClientType, c.Password, c.Email, c.Passport, c.CreditID,
	).Scan(&c.ID)
	return err
}

func (r *ClientRepository) Update(c *models.Client) error {
	_, err := db.DB.Exec(
		"UPDATE client SET phone=$1, login=$2, fullname=$3, clienttype=$4, password=$5, email=$6, passport=$7, creditid=$8 WHERE id=$9",
		c.Phone, c.Login, c.FullName, c.ClientType, c.Password, c.Email, c.Passport, c.CreditID, c.ID,
	)
	return err
}

func (r *ClientRepository) Delete(id int) error {
	_, err := db.DB.Exec("DELETE FROM client WHERE id = $1", id)
	return err
}
