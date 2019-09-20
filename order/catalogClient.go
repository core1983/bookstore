package order

import (
	"context"
	"google.golang.org/grpc"
	catalogpb "order/pb/pbcatalog"
)

type CatalogClient struct {
	conn    *grpc.ClientConn
	service catalogpb.CatalogServiceClient
}

func NewCatalogClient(url string) (*CatalogClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := catalogpb.NewCatalogServiceClient(conn)
	return &CatalogClient{conn, c}, nil
}

func (c *CatalogClient) Close() {
	c.conn.Close()
}

func (c *CatalogClient) CreateBook(ctx context.Context, isbn string, title string, author string, price float64) (*Book, error) {
	r, err := c.service.CreateBook(
		ctx,
		&catalogpb.CreateBookRequest{
			Isbn:   isbn,
			Title:  title,
			Author: author,
			Price:  price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Book{
		ID:     r.Book.Id,
		ISBN:   r.Book.Isbn,
		Title:  r.Book.Title,
		Author: r.Book.Author,
		Price:  r.Book.Price,
	}, nil
}

func (c *CatalogClient) GetBookByID(ctx context.Context, id string) *Book {
	r, err := c.service.GetBookByID(
		ctx,
		&catalogpb.GetBookRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil
	}

	return &Book{
		ID:     r.Book.Id,
		Title:  r.Book.Title,
		Author: r.Book.Author,
		Price:  r.Book.Price,
	}
}
