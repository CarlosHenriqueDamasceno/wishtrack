package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const QueryExecTimeout = time.Second * 10

type DatabaseRepository struct {
	connection *sql.DB
}

func NewDatabaseRepository(connection *sql.DB) *DatabaseRepository {
	return &DatabaseRepository{
		connection: connection,
	}
}

func (r *DatabaseRepository) Create(ctx context.Context, user *User) error {
	stm, err := r.connection.Prepare(`
	INSERT INTO users
		(id, name, email, password)
	VALUES
		(?, ?, ?, ?)
	RETURNING created_at`)

	if err != nil {
		return err
	}
	defer stm.Close()

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	err = stm.QueryRowContext(ctx, user.ID.String(), user.Name, user.Email, user.Password).Scan(&user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *DatabaseRepository) Find(ctx context.Context, id uuid.UUID) (*User, error) {
	stm, err := r.connection.Prepare("select id, name, email, password, created_at from users where id = ?")
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	row := stm.QueryRowContext(ctx, id.String())
	user := &User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

func (r *DatabaseRepository) IsEmailAlreadyTaken(ctx context.Context, email Email) (bool, error) {
	stm, err := r.connection.Prepare("select count(id) from users where email = ?")
	if err != nil {
		return false, err
	}
	defer stm.Close()

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	res, err := stm.ExecContext(ctx, email)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
