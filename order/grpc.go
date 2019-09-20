package order

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	orderpb "order/pb/pborder"
)

type grpcServer struct {
	service       Service
	userClient    *UserClient
	catalogClient *CatalogClient
}

func (gS *grpcServer) CreateOrder(ctx context.Context, r *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// Check if account exists
	_, err := gS.userClient.GetAccount(ctx, r.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get ordered products
	var bookIDs []string
	for _, b := range r.Books {
		bookIDs = append(bookIDs, b.BookId)
	}
	var orderedBooks []Book
	for _, id := range bookIDs {
		book := gS.catalogClient.GetBookByID(ctx, id)
		orderedBooks = append(orderedBooks, *book)
	}
	// Construct products
	var books []OrderedBook
	for _, p := range orderedBooks {
		orderedBook := OrderedBook{
			ID:       p.ID,
			Quantity: 0,
			Price:    p.Price,
			Title:    p.Title,
			Author:   p.Author,
		}
		for _, rp := range r.Books {
			if rp.BookId == p.ID {
				orderedBook.Quantity = rp.Quantity
				break
			}
		}

		if orderedBook.Quantity != 0 {
			books = append(books, orderedBook)
		}
	}

	// Call service implementation
	order, err := gS.service.CreateOrder(ctx, r.UserId, books)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderProto := &orderpb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		TotalPrice: order.TotalPrice,
		Books:      []*orderpb.Order_OrderBook{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Books {
		orderProto.Books = append(orderProto.Books, &orderpb.Order_OrderBook{
			Id:       p.ID,
			Isbn:     p.ISBN,
			Title:    p.Title,
			Author:   p.Author,
			Price:    p.Price,
			Quantity: p.Quantity,
		})
	}
	return &orderpb.CreateOrderResponse{
		Order: orderProto,
	}, nil
}

func (gS *grpcServer) GetOrdersForUser(ctx context.Context, r *orderpb.GetOrdersForUserRequest) (*orderpb.GetOrdersForUserResponse, error) {
	userOrders, err := gS.service.GetOrdersForUser(ctx, r.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get all ordered products
	orderIDMap := map[string]bool{}
	for _, o := range userOrders {
		for _, p := range o.Books {
			orderIDMap[p.ID] = true
		}
	}
	bookIDs := []string{}
	for id := range orderIDMap {
		bookIDs = append(bookIDs, id)
	}
	var orderedBooks []Book
	for _, id := range bookIDs {
		book := gS.catalogClient.GetBookByID(ctx, id)
		orderedBooks = append(orderedBooks, *book)
	}

	// Construct orders
	var orders []*orderpb.Order
	for _, o := range userOrders {
		// Encode order
		op := &orderpb.Order{
			UserId:     o.UserID,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Books:      []*orderpb.Order_OrderBook{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, book := range o.Books {
			// Populate product fields
			for _, p := range orderedBooks {
				if p.ID == book.ID {
					book.ISBN = p.ISBN
					book.Title = p.Title
					book.Author = p.Author
					book.Price = p.Price
					break
				}
			}

			op.Books = append(op.Books, &orderpb.Order_OrderBook{
				Id:       book.ID,
				Isbn:     book.ISBN,
				Title:    book.Title,
				Author:   book.Author,
				Price:    book.Price,
				Quantity: book.Quantity,
			})
		}

		orders = append(orders, op)
	}
	return &orderpb.GetOrdersForUserResponse{Orders: orders}, nil
}

func ListenGRPC(s Service, userServerURL, catalogServerURL string, port int) error {
	userClient, err := NewUserClient(userServerURL)
	if err != nil {
		return err
	}

	catalogClient, err := NewCatalogClient(catalogServerURL)
	if err != nil {
		userClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		userClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(serv, &grpcServer{
		s,
		userClient,
		catalogClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}
