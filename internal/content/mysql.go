package content

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
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
		return database.ParseDatabaseError(err)
	}

	persistedContent, err := r.Find(ctx, content.ID)
	if err != nil {
		return database.ParseDatabaseError(err)
	}

	content.CreatedAt = persistedContent.CreatedAt
	content.UpdatedAt = persistedContent.UpdatedAt

	return nil
}

func (r *DatabaseRepository) Find(ctx context.Context, id uuid.UUID) (*Content, error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, rate, comment, created_at, updated_at
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
		&content.Rate,
		&content.Comment,
		&content.CreatedAt,
		&content.UpdatedAt,
	)
	if err != nil {
		return nil, database.ParseDatabaseError(err)
	}

	content.Genres = strings.Split(genres, "|")

	return content, nil
}

func (r *DatabaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM contents WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(ctx, query, id)
	if err != nil {
		return database.ParseDatabaseError(err)
	}

	return nil
}

func (r *DatabaseRepository) Feed(ctx context.Context, userId uuid.UUID) ([]*Content, error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE user_id = ? AND rate IS NULL ORDER BY wish_level DESC LIMIT 5
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	rows, err := r.connection.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, database.ParseDatabaseError(err)
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
			return nil, database.ParseDatabaseError(err)
		}

		content.Genres = strings.Split(genres, "|")

		feed = append(feed, content)
	}

	return feed, nil
}

func (r *DatabaseRepository) List(ctx context.Context, userId uuid.UUID, pagination query.PaginationInput, filters ContentListFilters) (data []*Content, total uint64, err error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE user_id = ?
	`
	args := []any{userId}
	query, args = r.appendFilters(query, args, filters)
	offset := pagination.Limit * (pagination.Page - 1)
	args = append(args, pagination.Limit, offset)

	query += " LIMIT ? OFFSET ?"

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	rows, err := r.connection.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, database.ParseDatabaseError(err)
	}

	var list []*Content

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
			return nil, 0, database.ParseDatabaseError(err)
		}

		content.Genres = strings.Split(genres, "|")

		list = append(list, content)
	}

	query = `SELECT COUNT(id) as total FROM contents`
	err = r.connection.QueryRow(query).Scan(&total)
	if err != nil {
		return nil, 0, database.ParseDatabaseError(err)
	}

	return list, total, nil
}

func (r *DatabaseRepository) Update(ctx context.Context, content *Content) error {
	query := `UPDATE contents
				SET name = ?, category = ?, genres = ?, summary = ?, wish_level = ?, rate = ?, comment = ?, updated_at = ?
				WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	updatedAt := time.Now()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		content.Name,
		content.Category,
		strings.Join(content.Genres, "|"),
		content.Summary,
		content.WishLevel,
		content.Rate,
		content.Comment,
		updatedAt,
		content.ID,
	)
	if err != nil {
		return database.ParseDatabaseError(err)
	}

	content.UpdatedAt = updatedAt
	return nil
}

func (r *DatabaseRepository) appendFilters(query string, args []any, filters ContentListFilters) (string, []any) {
	if filters.Category != nil {
		query += " AND category = ?"
		args = append(args, filters.Category)
	}

	if filters.Watched != nil {
		if *filters.Watched {
			query += " AND rate IS NOT NULL"
		} else {
			query += " AND rate IS NULL"
		}
	}

	if filters.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filters.Name+"%")
	}

	if filters.Summary != nil {
		query += " AND summary LIKE ?"
		args = append(args, "%"+*filters.Summary+"%")
	}

	if filters.WishLevel != nil {
		query += " AND wish_level >= ?"
		args = append(args, *filters.WishLevel)
	}

	if filters.Genres != nil && len(*filters.Genres) > 0 {
		var conditions []string
		for _, genre := range *filters.Genres {
			conditions = append(conditions, fmt.Sprintf(
				"genres LIKE '%%|%s|%%' OR genres LIKE '%s|%%' OR genres LIKE '%%|%s' OR genres = '%s'",
				genre,
				genre,
				genre,
				genre,
			))
		}
		value := strings.Join(conditions, " OR ")
		query += " AND " + value
	}

	return query, args
}
