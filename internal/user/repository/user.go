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
	GetByEmail(email string) (model.User, error)
	UpdateEmail(id int64, email string) error
	IsEmailTaken(email string, excludeID int64) (bool, error)
	UpdatePassword(id int64, hashedPassword string) error


	
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll() ([]model.User, error) {
	rows, err := r.db.Query(`
		SELECT user_id, first_name, lastname, password, phone_number, email, photo, created_at, updated_at, role
		FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Password,
			&user.PhoneNumber, &user.Email, &user.Photo,
			&user.CreatedAt, &user.UpdatedAt, &user.Role,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetByID(id int64) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(`
		SELECT user_id, first_name, lastname, password, phone_number, email, photo, created_at, updated_at, role
		FROM users WHERE user_id = ?`, id).
		Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Password,
			&user.PhoneNumber, &user.Email, &user.Photo,
			&user.CreatedAt, &user.UpdatedAt, &user.Role,
		)
	return user, err
}

func (r *userRepo) Create(user model.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (first_name, email, password)
		VALUES (?, ?, ?)`,
		user.FirstName, user.Email, user.Password,
	)
	return err
}

func (r *userRepo) Update(u model.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET first_name=?, lastname=?, phone_number=?, photo=?, role=?, updated_at=NOW()
		WHERE user_id=?`,
		u.FirstName, u.LastName, u.PhoneNumber, u.Photo, u.Role, u.ID)
	return err
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE user_id=?", id)
	return err
}

func (r *userRepo) GetByEmail(email string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(`
	SELECT user_id, first_name, lastname, password, phone_number, email, photo, created_at, updated_at, role 
	FROM users WHERE TRIM(LOWER(email)) = TRIM(LOWER(?))`, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.PhoneNumber, &u.Email, &u.Photo, &u.CreatedAt, &u.UpdatedAt, &u.Role)

	return u, err
}

//  Update Email
func (r *userRepo) UpdateEmail(id int64, email string) error {
	_, err := r.db.Exec(`UPDATE users SET email = ?, updated_at = NOW() WHERE user_id = ?`, email, id)
	return err
}

//  Update Password 
func (r *userRepo) UpdatePassword(id int64, hashedPassword string) error {
	_, err := r.db.Exec(`UPDATE users SET password = ?, updated_at = NOW() WHERE user_id = ?`, hashedPassword, id)
	return err
}

// 
func (r *userRepo) IsEmailTaken(email string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE email = ? AND user_id != ?`, email, excludeID).Scan(&count)
	return count > 0, err
}
