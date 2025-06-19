package service

import (
	"context"
	"fmt"
	"strings"

	"news-service/internal/cache"
	"news-service/internal/domain"
	"news-service/internal/repository"
	"news-service/pkg/errors"
)

type NewsService struct {
	repo  repository.NewsRepository
	cache *cache.Cache
}

func NewNewsService(repo repository.NewsRepository, cache *cache.Cache) *NewsService {
	return &NewsService{
		repo:  repo,
		cache: cache,
	}
}

func (s *NewsService) CreateNews(ctx context.Context, slug, title, content string) (*domain.News, error) {
	// Валидация входных данных
	if err := s.validateNewsData(slug, title, content); err != nil {
		return nil, err
	}

	news := &domain.News{
		Slug:    slug,
		Title:   title,
		Content: content,
	}

	// Сохраняем в БД
	if err := s.repo.Create(ctx, news); err != nil {
		return nil, err
	}

	// Добавляем в кеш
	s.cache.Set(s.getCacheKey(slug), news)

	return news, nil
}

func (s *NewsService) GetNews(ctx context.Context, slug string) (*domain.News, error) {
	if slug == "" {
		return nil, errors.ErrInvalidSlug
	}

	// Сначала проверяем кеш
	cacheKey := s.getCacheKey(slug)
	if cached, exists := s.cache.Get(cacheKey); exists {
		if news, ok := cached.(*domain.News); ok {
			return news, nil
		}
	}

	// Если в кеше нет, запрашиваем из БД
	news, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Кешируем результат
	s.cache.Set(cacheKey, news)

	return news, nil
}

func (s *NewsService) GetNewsList(ctx context.Context, page, limit int) ([]*domain.News, int64, error) {
	// Валидация пагинации
	if page < 1 || limit < 1 || limit > 100 {
		return nil, 0, errors.ErrInvalidPagination
	}

	offset := (page - 1) * limit

	// Проверяем кеш для списка
	listCacheKey := s.getListCacheKey(page, limit)
	if cached, exists := s.cache.Get(listCacheKey); exists {
		if result, ok := cached.(ListCacheItem); ok {
			return result.News, result.Total, nil
		}
	}

	// Запрашиваем из БД
	newsList, total, err := s.repo.GetList(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Кешируем результат
	s.cache.Set(listCacheKey, ListCacheItem{
		News:  newsList,
		Total: total,
	})

	// Также кешируем индивидуальные новости
	for _, news := range newsList {
		s.cache.Set(s.getCacheKey(news.Slug), news)
	}

	return newsList, total, nil
}

func (s *NewsService) UpdateNews(ctx context.Context, slug, title, content string) (*domain.News, error) {
	// Валидация входных данных
	if err := s.validateNewsData(slug, title, content); err != nil {
		return nil, err
	}

	news := &domain.News{
		Title:   title,
		Content: content,
	}

	// Обновляем в БД
	if err := s.repo.Update(ctx, slug, news); err != nil {
		return nil, err
	}

	// Инвалидируем кеш для этой новости
	s.cache.Delete(s.getCacheKey(slug))

	// Инвалидируем кеш списков (упрощенный подход)
	s.invalidateListCache()

	return news, nil
}

func (s *NewsService) DeleteNews(ctx context.Context, slug string) error {
	if slug == "" {
		return errors.ErrInvalidSlug
	}

	// Удаляем из БД
	if err := s.repo.Delete(ctx, slug); err != nil {
		return err
	}

	// Удаляем из кеша
	s.cache.Delete(s.getCacheKey(slug))

	// Инвалидируем кеш списков
	s.invalidateListCache()

	return nil
}

func (s *NewsService) validateNewsData(slug, title, content string) error {
	if slug == "" || len(slug) > 255 {
		return errors.ErrInvalidSlug
	}

	if title == "" || len(title) > 500 {
		return errors.ErrInvalidTitle
	}

	if content == "" {
		return errors.ErrInvalidContent
	}

	// Простая валидация slug (только буквы, цифры, дефисы)
	slug = strings.ToLower(slug)
	for _, char := range slug {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return errors.ErrInvalidSlug
		}
	}

	return nil
}

func (s *NewsService) getCacheKey(slug string) string {
	return fmt.Sprintf("news:%s", slug)
}

func (s *NewsService) getListCacheKey(page, limit int) string {
	return fmt.Sprintf("news_list:%d:%d", page, limit)
}

func (s *NewsService) invalidateListCache() {
	// Простая реализация - можно улучшить, храня список ключей списков
	// Для простоты очистим весь кеш списков (в реальном проекте лучше более точечную инвалидацию)
}

type ListCacheItem struct {
	News  []*domain.News
	Total int64
}
