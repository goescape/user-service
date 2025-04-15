package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	fault "user-svc/helpers/fault"
	"user-svc/model"

	"github.com/google/uuid"
)

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *userStore {
	return &userStore{
		db: db,
	}
}

type UserRepository interface {
	Insert(user model.RegisterUser) (*uuid.UUID, error)
	Detail(req model.GetUserDetailRequest) (*model.User, error)
	ExistsByName(name string) (bool, error)
}

func (s *userStore) Insert(user model.RegisterUser) (*uuid.UUID, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fault.Custom(
			http.StatusConflict,
			fault.ErrConflict,
			fmt.Sprintf("failed start db transaction: %v", err.Error()))
	}
	defer tx.Rollback()

	baseQuery := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`

	var userId uuid.UUID
	if err := tx.QueryRow(baseQuery, user.Name, user.Email, user.Password).Scan(&userId); err != nil {
		tx.Rollback()
		return nil, fault.Custom(
			http.StatusConflict,
			fault.ErrConflict,
			fmt.Sprintf("failed to insert user: %v", err.Error()),
		)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fault.Custom(
			http.StatusConflict,
			fault.ErrConflict,
			fmt.Sprintf("failed to commit transaction: %v", err),
		)
	}

	return &userId, nil
}

func (s *userStore) Detail(req model.GetUserDetailRequest) (*model.User, error) {
	baseQuery := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE `
	var args []interface{}
	var conditions []string
	argPos := 1

	if req.UserId != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argPos))
		args = append(args, req.UserId)
		argPos++
	}

	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", argPos))
		args = append(args, req.Name)
		argPos++
	}

	if req.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argPos))
		args = append(args, req.Email)
		argPos++
	}

	if len(conditions) == 0 {
		return nil, fault.Custom(
			http.StatusBadRequest,
			fault.ErrBadRequest,
			"at least one filter (user_id, name, or email) must be provided",
		)
	}

	query := baseQuery + strings.Join(conditions, " AND ")

	var user model.User
	err := s.db.QueryRow(query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fault.Custom(
				http.StatusNotFound,
				fault.ErrNotFound,
				"user not found based on provided filters",
			)
		}
		return nil, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to get user detail: %v", err),
		)
	}

	return &user, nil
}

func (s *userStore) ExistsByName(name string) (bool, error) {
	baseQuery := `SELECT COUNT(*) FROM users WHERE name = $1`

	var count int
	err := s.db.QueryRow(baseQuery, name).Scan(&count)
	if err != nil {
		return false, fault.Custom(
			http.StatusInternalServerError,
			fault.ErrInternalServer,
			fmt.Sprintf("failed to count users by name '%s': %v", name, err),
		)
	}

	return count > 0, nil
}
