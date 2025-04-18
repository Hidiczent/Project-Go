package repository

import (
	"database/sql"
	"myapp/internal/user/model"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByID(id int64) (model.User, error)
	Create(user model.User) error
	Update(user model.User) error
	Delete(id int64) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}
func (r *userRepo) GetAll() ([]model.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, password FROM users") // ✅ เพิ่ม password
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password) // ✅ ต้องตรงกับ SELECT
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetByID(id int64) (model.User, error) {
	var user model.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password) // ✅ เพิ่ม password
	return user, err
}


func (r *userRepo) Create(user model.User) error {
	_, err := r.db.Exec("INSERT INTO users(name, email, password) VALUES(?, ?, ?)", user.Name, user.Email, user.Password)
	return err
}

func (r *userRepo) Update(user model.User) error {
	_, err := r.db.Exec("UPDATE users SET name=?, email=?, password=? WHERE id=?", user.Name, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=?", id)
	return err
}
