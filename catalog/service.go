package catalog

import (
	"context"
	"log"
)
import "github.com/rs/xid"

type Service interface {
	CreateBook(ctx context.Context, isbn string, title string, author string, price float64) (*Book, error)
	GetBookByID(ctx context.Context, id string) (*Book, error)
	GetBooksList(ctx context.Context) ([]Book, error)
}

type bookService struct {
	bookRepository Repository
}

func NewBookService(r Repository) Service {
	return &bookService{r}
}

func (bS *bookService) CreateBook(ctx context.Context, isbn string, title string, author string, price float64) (*Book, error) {
	newBook := &Book{
		ID:     xid.New().String(),
		ISBN:   isbn,
		Title:  title,
		Author: author,
		Price:  price,
	}
	log.Print(newBook)
	if err := bS.bookRepository.CreateBook(ctx, *newBook); err != nil {
		return nil, err
	}
	return newBook, nil
}

func (bS *bookService) GetBookByID(ctx context.Context, id string) (*Book, error) {
	return bS.bookRepository.GetBookByID(ctx, id)
}

func (bS *bookService) GetBooksList(ctx context.Context) ([]Book, error) {
	return bS.bookRepository.BooksList(ctx)
}
