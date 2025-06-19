package repository

import (
	"context"
	"news-service/internal/domain"
)

type NewsRepository interface {
	Create(ctx context.Context, news *domain.News) error
	GetBySlug(ctx context.Context, slug string) (*domain.News, error)
	GetList(ctx context.Context, offset, limit int) ([]*domain.News, int64, error)
	Update(ctx context.Context, slug string, news *domain.News) error
	Delete(ctx context.Context, slug string) error
}
