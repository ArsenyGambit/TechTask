package main

import (
	"context"
	"log"
	"time"

	"news-service/internal/config"
	"news-service/internal/domain"
	"news-service/internal/repository/postgres"
	"news-service/pkg/database"
)

func main() {
	cfg, err := config.LoadDefault()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := postgres.NewNewsRepository(db)
	ctx := context.Background()

	// Тестовые данные
	testNews := []*domain.News{
		{
			Slug:    "first-news",
			Title:   "Первая новость",
			Content: "Это содержимое первой новости. Здесь много интересной информации.",
		},
		{
			Slug:    "breaking-news",
			Title:   "Срочные новости",
			Content: "Важная новость, которую все должны знать. Подробности внутри.",
		},
		{
			Slug:    "tech-update",
			Title:   "Обновление технологий",
			Content: "Новые технологии изменяют мир. В этой статье рассказываем о последних трендах.",
		},
		{
			Slug:    "sports-news",
			Title:   "Спортивные новости",
			Content: "Результаты последних спортивных событий и анонсы предстоящих матчей.",
		},
		{
			Slug:    "weather-forecast",
			Title:   "Прогноз погоды",
			Content: "Погода на завтра и ближайшие дни. Не забудьте взять зонт!",
		},
	}

	log.Println("Начинаем заполнение БД тестовыми данными...")

	for _, news := range testNews {
		if err := repo.Create(ctx, news); err != nil {
			log.Printf("Ошибка при создании новости %s: %v", news.Slug, err)
			continue
		}
		log.Printf("Создана новость: %s", news.Title)

		// Небольшая задержка для разных created_at
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Заполнение БД завершено!")
}
