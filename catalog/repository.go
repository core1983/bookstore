package catalog

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Repository interface {
	CreateBook(ctx context.Context, b Book) error
	GetBookByID(ctx context.Context, id string) (*Book, error)
	BooksList(ctx context.Context) ([]Book, error)
	Close()
}

type bookRepository struct {
	db *sql.DB
}

func (bR *bookRepository) CreateBook(ctx context.Context, b Book) error {
	_, err := bR.db.ExecContext(ctx, "INSERT INTO books(id,isbn, title, author, price ) VALUES($1,$2,$3,$4,$5)", b.ID, b.ISBN, b.Title, b.Author, b.Price)
	return err
}

func (bR *bookRepository) GetBookByID(ctx context.Context, id string) (*Book, error) {
	row := bR.db.QueryRowContext(ctx, "SELECT id,isbn, title, author, price FROM books WHERE id = $1", id)
	b := &Book{}
	if err := row.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Price); err != nil {
		log.Printf("err %v", err)
		return nil, err
	}
	log.Print(b)
	return b, nil
}

func (bR *bookRepository) BooksList(ctx context.Context) ([]Book, error) {
	rows, err := bR.db.QueryContext(ctx, "SELECT id,isbn, title, author, price FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		b := &Book{}
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Price); err == nil {
			books = append(books, *b)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return books, nil
}

func (bR *bookRepository) Close() {
	bR.db.Close()
}

func NewBookRepository(connectionString string) (Repository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	log.Print(db.Stats())
	return &bookRepository{db}, nil
}
