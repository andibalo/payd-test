package repository

import (
	"database/sql"
	"github.com/andibalo/payd-test/backend/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

type GetUserListFilter struct {
	Role string
}

func (r *userRepository) Save(user *model.User) error {
	query := `
		INSERT INTO users (
			first_name, last_name, email, password, role, created_by
		)
		VALUES ( ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, first_name, last_name, email, password, role, created_by, created_at, 
		       updated_by, updated_at
		FROM users 
		WHERE email = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role,
		&user.CreatedBy, &user.CreatedAt, &user.UpdatedBy, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetList(filter GetUserListFilter) ([]model.User, error) {
	users := []model.User{}

	query := `
		SELECT id, first_name, last_name, email, password, role, created_by, created_at, 
			   updated_by, updated_at 
		FROM users WHERE deleted_at IS NULL
	`

	var args []interface{}

	if filter.Role != "" {
		query += " AND role = ?"
		args = append(args, filter.Role)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role,
			&user.CreatedBy, &user.CreatedAt, &user.UpdatedBy, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
