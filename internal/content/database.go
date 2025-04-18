package content

import (
	"context"
	"database/sql"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
			($1, $2, $3, $4, $5, $6, $7)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		content.ID,
		content.Name,
		content.Category,
		pq.Array(content.Genres),
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
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	content := &Content{}
	var genres []string
	err := r.connection.QueryRowContext(ctx, query, id).Scan(
		&content.ID,
		&content.Name,
		&content.Category,
		pq.Array(&genres),
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

	content.Genres = genres

	return content, nil
}

func (r *DatabaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM contents WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	_, err := r.connection.ExecContext(ctx, query, id)
	if err != nil {
		return database.ParseDatabaseError(err)
	}

	return nil
}

func (r *DatabaseRepository) Suggestions(ctx context.Context, userId uuid.UUID, limit int) ([]*Content, error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE user_id = $1 AND rate IS NULL ORDER BY wish_level DESC LIMIT $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	rows, err := r.connection.QueryContext(ctx, query, userId, limit)
	if err != nil {
		return nil, database.ParseDatabaseError(err)
	}

	var suggestions []*Content

	for rows.Next() {
		content := &Content{}
		var genres []string
		err := rows.Scan(
			&content.ID,
			&content.Name,
			&content.Category,
			pq.Array(&genres),
			&content.Summary,
			&content.WishLevel,
			&content.UserID,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, database.ParseDatabaseError(err)
		}

		content.Genres = genres
		suggestions = append(suggestions, content)
	}

	return suggestions, nil
}

func (r *DatabaseRepository) List(ctx context.Context, userId uuid.UUID, pagination query.PaginationInput, filters ContentListFilters) (data []*Content, total uint64, err error) {
	query := `
		SELECT
			id, name, category, genres, summary, wish_level, user_id, created_at, updated_at
		FROM contents
		WHERE user_id = $1
		AND (category ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR summary ILIKE '%' || $2 || '%' OR $2 IS NULL)
		AND (genres && $3 OR $3 IS NULL)
		AND (wish_level >= $4 OR $4 IS NULL)
		AND (
			($5 IS TRUE AND rate IS NOT NULL) OR
			($5 IS FALSE AND rate IS NULL) OR
			($5 IS NULL)
		)
		ORDER BY created_at DESC
		LIMIT $6 OFFSET $7
	`

	offset := pagination.Limit * (pagination.Page - 1)

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	rows, err := r.connection.QueryContext(
		ctx,
		query,
		userId,
		filters.Search,
		pq.Array(filters.Genres),
		filters.WishLevel,
		filters.Watched,
		pagination.Limit,
		offset,
	)
	if err != nil {
		return nil, 0, database.ParseDatabaseError(err)
	}

	var list []*Content

	for rows.Next() {
		content := &Content{}
		var genres []string
		err := rows.Scan(
			&content.ID,
			&content.Name,
			&content.Category,
			pq.Array(&genres),
			&content.Summary,
			&content.WishLevel,
			&content.UserID,
			&content.CreatedAt,
			&content.UpdatedAt,
		)
		if err != nil {
			return nil, 0, database.ParseDatabaseError(err)
		}

		content.Genres = genres
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
				SET name = $1, category = $2, genres = $3, summary = $4, wish_level = $5, rate = $6, comment = $7, updated_at = $8
				WHERE id = $9`

	ctx, cancel := context.WithTimeout(ctx, QueryExecTimeout)
	defer cancel()

	updatedAt := time.Now()

	_, err := r.connection.ExecContext(
		ctx,
		query,
		content.Name,
		content.Category,
		pq.Array(content.Genres),
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
