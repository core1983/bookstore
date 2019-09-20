package order

import (
	"context"
	"google.golang.org/grpc"
	userpb "order/pb/pbuser"
)

type UserClient struct {
	conn    *grpc.ClientConn
	service userpb.UserServiceClient
}

func NewUserClient(url string) (*UserClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := userpb.NewUserServiceClient(conn)
	return &UserClient{conn, c}, nil
}

func (c *UserClient) Close() {
	c.conn.Close()
}

func (c *UserClient) CreateUser(ctx context.Context, name string) (*User, error) {
	r, err := c.service.CreateUser(
		ctx,
		&userpb.CreateUserRequest{UserName: name},
	)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       r.User.Id,
		UserName: r.User.UserName,
	}, nil
}

func (c *UserClient) GetAccount(ctx context.Context, id string) (*User, error) {
	r, err := c.service.GetUserByID(
		ctx,
		&userpb.GetUserRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       r.User.Id,
		UserName: r.User.UserName,
	}, nil
}

func (c *UserClient) GetAccounts(ctx context.Context) ([]User, error) {
	r, err := c.service.GetUsersList(ctx, &userpb.GetUsersListRequest{})
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, a := range r.Users {
		users = append(users, User{
			ID:       a.Id,
			UserName: a.UserName,
		})
	}
	return users, nil
}
