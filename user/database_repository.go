package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type DatabaseRepository struct {
	connection *sql.DB
}

func NewDatabaseRepository(connection *sql.DB) *DatabaseRepository {
	return &DatabaseRepository{
		connection: connection,
	}
}

func (r *DatabaseRepository) Create(user *User) error {
	stm, err := r.connection.Prepare("insert into users (id, name, email, password, created_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(user.ID.String(), user.Name, user.Email, user.password, user.CreatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func (r *DatabaseRepository) Find(id uuid.UUID) (*User, error) {
	stm, err := r.connection.Prepare("select id, name, email, password, created_at from users where id = ?")
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	row := stm.QueryRow(id.String())
	user := &User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.password, &user.CreatedAt)
	return user, err
}
