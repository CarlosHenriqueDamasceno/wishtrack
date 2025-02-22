package content

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
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

func (r *DatabaseRepository) Create(ctx context.Context, content *Content) error {
	query := `
		INSERT INTO contents
			(id, name, category, genres, summary, wish_level, user_id)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		content.ID,
		content.Name,
		content.Category,
		strings.Join(content.Genres, "|"),
		content.Summary,
		content.WishLevel,
		content.UserID,
	)
	if err != nil {
		return database.WrapMysqlError(err)
	}

	persistedContent, err := r.Find(ctx, content.ID)
	if err != nil {
		return database.WrapMysqlError(err)
	}

	content.CreatedAt = persistedContent.CreatedAt
	content.UpdatedAt = persistedContent.UpdatedAt

	return nil
}

func (r *DatabaseRepository) Find(ctx context.Context, id uuid.UUID) (*Content, error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	var genres string

	content := &Content{}
	err := r.connection.QueryRowContext(ctx, query, id).Scan(
		&content.ID,
		&content.Name,
		&content.Category,
		&genres,
		&content.Summary,
		&content.WishLevel,
		&content.UserID,
		&content.CreatedAt,
		&content.UpdatedAt,
	)
	if err != nil {
		return nil, database.WrapMysqlError(err)
	}

	content.Genres = strings.Split(genres, "|")

	return content, nil
}

func (r *DatabaseRepository) Feed(ctx context.Context, userId uuid.UUID) ([]*Content, error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE user_id = ? ORDER BY wish_level DESC LIMIT 5
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	rows, err := r.connection.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, database.WrapMysqlError(err)
	}

	var feed []*Content

	for rows.Next() {
		var genres string
		content := &Content{}

		err := rows.Scan(
			&content.ID,
			&content.Name,
			&content.Category,
			&genres,
			&content.Summary,
			&content.WishLevel,
			&content.UserID,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, database.WrapMysqlError(err)
		}

		content.Genres = strings.Split(genres, "|")

		feed = append(feed, content)
	}

	return feed, nil
}
