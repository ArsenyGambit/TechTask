package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"news-service/internal/domain"
	"news-service/internal/service"
	"news-service/pkg/errors"
	pb "news-service/proto/news"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedNewsServiceServer // Добавили встраивание
	newsService                       *service.NewsService
	grpcServer                        *grpc.Server
}

func NewServer(newsService *service.NewsService) *Server {
	return &Server{
		newsService: newsService,
		grpcServer:  grpc.NewServer(),
	}
}

func (s *Server) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// ГЛАВНОЕ ИЗМЕНЕНИЕ: Регистрируем сервис!
	pb.RegisterNewsServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer) // Для удобства тестирования

	log.Printf("gRPC server starting on %s", address)
	log.Printf("NewsService registered successfully")

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

// Изменяем сигнатуры методов на protobuf типы

func (s *Server) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) (*pb.CreateNewsResponse, error) {
	news, err := s.newsService.CreateNews(ctx, req.Slug, req.Title, req.Content)
	if err != nil {
		return &pb.CreateNewsResponse{
			Error: s.handleError(err),
		}, nil
	}

	return &pb.CreateNewsResponse{
		News: s.domainToProto(news),
	}, nil
}

func (s *Server) GetNews(ctx context.Context, req *pb.GetNewsRequest) (*pb.GetNewsResponse, error) {
	news, err := s.newsService.GetNews(ctx, req.Slug)
	if err != nil {
		return &pb.GetNewsResponse{
			Error: s.handleError(err),
		}, nil
	}

	return &pb.GetNewsResponse{
		News: s.domainToProto(news),
	}, nil
}

func (s *Server) GetNewsList(ctx context.Context, req *pb.GetNewsListRequest) (*pb.GetNewsListResponse, error) {
	newsList, total, err := s.newsService.GetNewsList(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return &pb.GetNewsListResponse{
			Error: s.handleError(err),
		}, nil
	}

	protoNews := make([]*pb.News, len(newsList))
	for i, news := range newsList {
		protoNews[i] = s.domainToProto(news)
	}

	return &pb.GetNewsListResponse{
		News:  protoNews,
		Total: total,
	}, nil
}

func (s *Server) UpdateNews(ctx context.Context, req *pb.UpdateNewsRequest) (*pb.UpdateNewsResponse, error) {
	news, err := s.newsService.UpdateNews(ctx, req.Slug, req.Title, req.Content)
	if err != nil {
		return &pb.UpdateNewsResponse{
			Error: s.handleError(err),
		}, nil
	}

	return &pb.UpdateNewsResponse{
		News: s.domainToProto(news),
	}, nil
}

func (s *Server) DeleteNews(ctx context.Context, req *pb.DeleteNewsRequest) (*pb.DeleteNewsResponse, error) {
	err := s.newsService.DeleteNews(ctx, req.Slug)
	if err != nil {
		return &pb.DeleteNewsResponse{
			Success: false,
			Error:   s.handleError(err),
		}, nil
	}

	return &pb.DeleteNewsResponse{
		Success: true,
	}, nil
}

// Изменяем только возвращаемый тип
func (s *Server) domainToProto(news *domain.News) *pb.News {
	if news == nil {
		return nil
	}

	return &pb.News{
		Slug:      news.Slug,
		Title:     news.Title,
		Content:   news.Content,
		CreatedAt: news.CreatedAt.Unix(),
		UpdatedAt: news.UpdatedAt.Unix(),
	}
}

func (s *Server) handleError(err error) string {
	switch err {
	case errors.ErrNewsNotFound:
		return "News not found"
	case errors.ErrDuplicateSlug:
		return "News with this slug already exists"
	case errors.ErrInvalidSlug:
		return "Invalid slug format"
	case errors.ErrInvalidTitle:
		return "Invalid title"
	case errors.ErrInvalidContent:
		return "Invalid content"
	case errors.ErrInvalidPagination:
		return "Invalid pagination parameters"
	default:
		log.Printf("Unexpected error: %v", err)
		return "Internal server error"
	}
}
