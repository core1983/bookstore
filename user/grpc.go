package user

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"user/pb"
)

type grpcServer struct {
	service Service
}

func (gS *grpcServer) CreateUser(ctx context.Context, r *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	u, err := gS.service.CreateUser(ctx, r.UserName)
	if err != nil {
		return nil, err
	}
	userProto := &userpb.User{
		Id:       u.ID,
		UserName: u.UserName,
	}
	userProto.CreatedAt, _ = u.CreatedAt.MarshalBinary()

	return &userpb.CreateUserResponse{
		User: userProto,
	}, nil
}

func (gS *grpcServer) GetUserByID(ctx context.Context, r *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	u, err := gS.service.GetUserByID(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	userProto := &userpb.User{
		Id:       u.ID,
		UserName: u.UserName,
	}
	userProto.CreatedAt, _ = u.CreatedAt.MarshalBinary()

	return &userpb.GetUserResponse{
		User: userProto,
	}, nil
}

func (gS *grpcServer) GetUsersList(ctx context.Context, r *userpb.GetUsersListRequest) (*userpb.GetUsersListResponse, error) {
	res, err := gS.service.GetUsersList(ctx)
	if err != nil {
		return nil, err
	}
	var users []*userpb.User
	for _, user := range res {
		userProto := &userpb.User{
			Id:       user.ID,
			UserName: user.UserName,
		}
		userProto.CreatedAt, _ = user.CreatedAt.MarshalBinary()
		users = append(users, userProto)
	}

	return &userpb.GetUsersListResponse{Users: users,}, nil
}

func NewGrpcServer(s Service, port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	userpb.RegisterUserServiceServer(gs, &grpcServer{s})
	reflection.Register(gs)
	return gs.Serve(list)
}
