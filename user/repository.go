package user

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Repository interface {
	CreateUser(ctx context.Context, u User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	UsersList(ctx context.Context) ([]User, error)
	Close()
}

type userRepository struct {
	db *sql.DB
}

func (uR *userRepository) CreateUser(ctx context.Context, u User) error {
	_, err := uR.db.ExecContext(ctx, "INSERT INTO users(id,user_name, created_at) VALUES($1,$2,$3)", u.ID, u.UserName, u.CreatedAt)
	log.Printf("fail %v", err)
	return err
}

func (uR *userRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := uR.db.QueryRowContext(ctx, "SELECT id, user_name,created_at FROM users WHERE id = $1", id)
	u := &User{}
	if err := row.Scan(&u.ID, &u.UserName, &u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (uR *userRepository) UsersList(ctx context.Context) ([]User, error) {
	rows, err := uR.db.QueryContext(ctx, "SELECT id, user_name, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		u := &User{}
		if err = rows.Scan(&u.ID, &u.UserName, &u.CreatedAt); err == nil {
			users = append(users, *u)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (uR *userRepository) Close() {
	uR.db.Close()
}

func NewUserRepository(connectionString string) (Repository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &userRepository{db}, nil
}
