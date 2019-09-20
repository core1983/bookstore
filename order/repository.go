package order

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
)

type Repository interface {
	CreateOrder(ctx context.Context, o Order) error
	GetOrdersForUser(ctx context.Context, userID string) ([]Order, error)
	Close()
}

type orderRepository struct {
	db *sql.DB
}

func (oR *orderRepository) CreateOrder(ctx context.Context, o Order) error {
	tx, err := oR.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO orders(id, created_at, user_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.UserID,
		o.TotalPrice, )

	if err != nil {
		return nil
	}
	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_books", "order_id", "book_id", "quantity"))
	for _, b := range o.Books {
		_, err = stmt.ExecContext(ctx, o.ID, b.ID, b.Quantity)
		if err != nil {
			return nil
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return nil
	}
	stmt.Close()

	return nil
}

func (oR *orderRepository) GetOrdersForUser(ctx context.Context, userID string) ([]Order, error) {
	rows, err := oR.db.QueryContext(ctx,
		`SELECT
      ord.id,
      ord.created_at,
      ord.user_id,
      ord.total_price::money::numeric::float8,
      op.book_id,
      op.quantity
    FROM orders ord JOIN order_books ob ON (o.id = ob.order_id)
    WHERE ord.user_id = $1
    ORDER BY ord.id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	order := &Order{}
	lastOrder := &Order{}
	orderedBook := &OrderedBook{}
	var books []OrderedBook

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.UserID,
			&order.TotalPrice,
			&orderedBook.ID,
			&orderedBook.Quantity,
		); err != nil {
			return nil, err
		}
		// Scan order
		if lastOrder.ID != "" && lastOrder.ID != order.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				UserID:     lastOrder.UserID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Books:      books,
			}
			orders = append(orders, newOrder)
			books = []OrderedBook{}
		}
		
		books = append(books, OrderedBook{
			ID:       orderedBook.ID,
			Quantity: orderedBook.Quantity,
		})

		*lastOrder = *order
	}

	if lastOrder != nil {
		newOrder := Order{
			ID:         lastOrder.ID,
			UserID:     lastOrder.UserID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Books:      books,
		}
		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (oR *orderRepository) Close() {
	oR.db.Close()
}

func NewOrderRepository(connectionString string) (Repository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &orderRepository{db}, nil
}
