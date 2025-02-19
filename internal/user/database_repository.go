package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const QueryExecTimeout = time.Second * 10

var ErrDuplicateEmail = errors.New("e-mail already in use")

type DatabaseRepository struct {
	connection *sql.DB
}

func NewDatabaseRepository(connection *sql.DB) *DatabaseRepository {
	return &DatabaseRepository{
		connection: connection,
	}
}

func (r *DatabaseRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users
			(id, name, email, password)
		VALUES
			(?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		user.ID.String(),
		user.Name,
		user.Email,
		user.password.value,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return ErrDuplicateEmail
		}
		return wrapMysqlError(err)
	}

	persistedUser, err := r.Find(ctx, user.ID)
	if err != nil {
		return wrapMysqlError(err)
	}

	user.CreatedAt = persistedUser.CreatedAt
	user.UpdatedAt = persistedUser.UpdatedAt

	return nil
}

func (r *DatabaseRepository) Find(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT
			id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	user := &User{}
	err := r.connection.QueryRowContext(ctx, query, id.String()).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.password.value,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, wrapMysqlError(err)
	}

	return user, nil
}

func (r *DatabaseRepository) IsEmailAlreadyTaken(ctx context.Context, email Email) (bool, error) {
	query := "SELECT COUNT(id) FROM users WHERE email = ?"

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	count := 0
	err := r.connection.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, wrapMysqlError(err)
	}

	return count > 0, nil
}

func wrapMysqlError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return errors.New(mysqlErr.Message)
	}
	return err
}
