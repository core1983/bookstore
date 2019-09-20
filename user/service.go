package user

import (
	"context"
	"github.com/rs/xid"
	"log"
	"time"
)

type Service interface {
	CreateUser(ctx context.Context, name string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUsersList(ctx context.Context) ([]User, error)
}

type userService struct {
	userRepository Repository
}

func NewUserService(r Repository) Service {
	return &userService{r}
}

func (uS *userService) CreateUser(ctx context.Context, name string) (*User, error) {
	newUser := &User{
		ID:        xid.New().String(),
		UserName:  name,
		CreatedAt: time.Now().UTC(),
	}
	if err := uS.userRepository.CreateUser(ctx, *newUser); err != nil {
		log.Fatal("not good")
		return nil, err
	}
	return newUser, nil
}

func (uS *userService) GetUserByID(ctx context.Context, id string) (*User, error) {
	return uS.userRepository.GetUserByID(ctx, id)
}

func (uS *userService) GetUsersList(ctx context.Context) ([]User, error) {
	return uS.userRepository.UsersList(ctx)
}
