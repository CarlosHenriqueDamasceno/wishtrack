package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
			($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.password.value,
	)
	if err != nil {
		log.Println(err)
		if postgresError, ok := err.(*pq.Error); ok && postgresError.Code == "23505" {
			return ErrDuplicateEmail
		}
		return database.ParseDatabaseError(err)
	}

	persistedUser, err := r.Find(ctx, user.ID)
	if err != nil {
		return database.ParseDatabaseError(err)
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
		WHERE id = $1
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
		log.Println(err)
		return nil, database.ParseDatabaseError(err)
	}

	return user, nil
}

func (r *DatabaseRepository) IsEmailAlreadyTaken(ctx context.Context, email Email) (bool, error) {
	query := "SELECT COUNT(id) FROM users WHERE email = $1"

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	count := 0
	err := r.connection.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, database.ParseDatabaseError(err)
	}

	return count > 0, nil
}

func (r *DatabaseRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
				SELECT
					id, name, email, password, created_at, updated_at
				FROM users
				WHERE email = $1
			`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	user := &User{}
	err := r.connection.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.password.value,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, database.ParseDatabaseError(err)
		}
	}

	return user, nil
}
