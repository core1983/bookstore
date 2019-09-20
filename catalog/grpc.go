package catalog

import (
	catalogpb "catalog/pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	service Service
}

func (gS *grpcServer) CreateBook(ctx context.Context, r *catalogpb.CreateBookRequest) (*catalogpb.CreateBookResponse, error) {
	b, err := gS.service.CreateBook(ctx, r.Isbn, r.Title, r.Author, r.Price)
	if err != nil {
		return nil, err
	}
	return &catalogpb.CreateBookResponse{
		Book: &catalogpb.Book{
			Id:     b.ID,
			Isbn:   b.ISBN,
			Title:  b.Title,
			Author: b.Author,
			Price:  b.Price,
		},
	}, nil
}

func (gS *grpcServer) GetBookByID(ctx context.Context, r *catalogpb.GetBookRequest) (*catalogpb.GetBookResponse, error) {
	b, err := gS.service.GetBookByID(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &catalogpb.GetBookResponse{
		Book: &catalogpb.Book{
			Id:     b.ID,
			Isbn:   b.ISBN,
			Title:  b.Title,
			Author: b.Author,
			Price:  b.Price,
		},
	}, nil
}

func (gS *grpcServer) GetBooksList(ctx context.Context, r *catalogpb.GetBooksRequest) (*catalogpb.GetBooksResponse, error) {
	res, err := gS.service.GetBooksList(ctx)
	if err != nil {
		return nil, err
	}
	var books []*catalogpb.Book
	for _, book := range res {
		books = append(books, &catalogpb.Book{
			Id:     book.ID,
			Isbn:   book.ISBN,
			Title:  book.Title,
			Author: book.Author,
			Price:  book.Price,
		}, )
	}
	return &catalogpb.GetBooksResponse{
		Books: books,
	}, nil
}

func NewGrpcServer(s Service, port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	catalogpb.RegisterCatalogServiceServer(gs, &grpcServer{s})
	reflection.Register(gs)
	return gs.Serve(list)
}
