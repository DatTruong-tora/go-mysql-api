package repository

import (
	"database/sql"
	"go-mysql-api/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, role, admin_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User = []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.AdminID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Create(u *models.User) error {
	query := "INSERT INTO users (name, email, role, admin_id) VALUES (?, ?, ?, ?)"
	res, err := r.DB.Exec(query, u.Name, u.Email, u.Role, u.AdminID)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var u models.User
	query := "SELECT id, name, email, role, admin_id FROM users WHERE id = ?"
	err := r.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.AdminID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(id int, u *models.User) error {
	query := "UPDATE users SET name=?, email=?, role=?, admin_id=? WHERE id=?"
	_, err := r.DB.Exec(query, u.Name, u.Email, u.Role, u.AdminID, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByAdminID(id int) error {
	rows, err := r.DB.Query("SELECT id, name, email, role, admin_id FROM users WHERE admin_id = ?", id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []models.User = []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.AdminID); err != nil {
			return err
		}
		users = append(users, u)
	}
	return nil
}
