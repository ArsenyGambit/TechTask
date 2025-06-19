package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"news-service/internal/domain"
	"news-service/internal/repository"
	"news-service/pkg/errors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) repository.NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Create(ctx context.Context, news *domain.News) error {
	query := `
		INSERT INTO news (slug, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	news.CreatedAt = now
	news.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query, news.Slug, news.Title, news.Content, news.CreatedAt, news.UpdatedAt)
	if err != nil {
		// Проверяем на дубликат по первичному ключу
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.ErrDuplicateSlug
		}
		return fmt.Errorf("failed to create news: %w", err)
	}

	return nil
}

func (r *newsRepository) GetBySlug(ctx context.Context, slug string) (*domain.News, error) {
	query := `
		SELECT slug, title, content, created_at, updated_at
		FROM news
		WHERE slug = $1
	`

	news := &domain.News{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&news.Slug,
		&news.Title,
		&news.Content,
		&news.CreatedAt,
		&news.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNewsNotFound
		}
		return nil, fmt.Errorf("failed to get news by slug: %w", err)
	}

	return news, nil
}

func (r *newsRepository) GetList(ctx context.Context, offset, limit int) ([]*domain.News, int64, error) {
	// Получаем общее количество записей
	countQuery := `SELECT COUNT(*) FROM news`
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get news count: %w", err)
	}

	// Получаем записи с пагинацией
	query := `
		SELECT slug, title, content, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get news list: %w", err)
	}
	defer rows.Close()

	var newsList []*domain.News
	for rows.Next() {
		news := &domain.News{}
		err := rows.Scan(
			&news.Slug,
			&news.Title,
			&news.Content,
			&news.CreatedAt,
			&news.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan news: %w", err)
		}
		newsList = append(newsList, news)
	}

	return newsList, total, nil
}

func (r *newsRepository) Update(ctx context.Context, slug string, news *domain.News) error {
	query := `
		UPDATE news
		SET title = $2, content = $3, updated_at = $4
		WHERE slug = $1
	`

	news.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query, slug, news.Title, news.Content, news.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update news: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrNewsNotFound
	}

	news.Slug = slug
	return nil
}

func (r *newsRepository) Delete(ctx context.Context, slug string) error {
	query := `DELETE FROM news WHERE slug = $1`

	result, err := r.db.ExecContext(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrNewsNotFound
	}

	return nil
}
