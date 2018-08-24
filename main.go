package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/ktr0731/evans-demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	id                  string
	firstName, lastName string
	gender              api.Gender
}

type UserService struct {
	store map[string]*User
	h     hash.Hash
}

func (s *UserService) CreateUsers(ctx context.Context, req *api.CreateUsersRequest) (*api.CreateUsersResponse, error) {
	for _, user := range req.GetUsers() {
		if _, err := io.WriteString(s.h, user.GetFirstName()+user.GetLastName()); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate id: %s", err)
		}
		id := fmt.Sprintf("%x", s.h.Sum(nil))
		s.store[id] = &User{
			id:        id,
			firstName: user.GetFirstName(),
			lastName:  user.GetLastName(),
			gender:    user.GetGender(),
		}
		s.h.Reset()
	}
	return &api.CreateUsersResponse{
		Message: "registration successful",
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *api.ListUsersRequest) (*api.ListUsersResponse, error) {
	users := make([]*api.ListUsersResponse_User, 0, len(s.store))
	for _, user := range s.store {
		users = append(users, &api.ListUsersResponse_User{
			Id:   user.id,
			Name: user.firstName + " " + user.lastName,
		})
	}
	return &api.ListUsersResponse{
		Users: users,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	user, ok := s.store[req.GetId()]
	if !ok {
		return nil, errors.New("no such user")
	}
	return &api.GetUserResponse{
		User: &api.User{
			FirstName: user.firstName,
			LastName:  user.lastName,
			Gender:    user.gender,
		},
	}, nil
}

func main() {
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	api.RegisterUserServiceServer(server, &UserService{
		store: map[string]*User{},
		h:     sha256.New(),
	})
	if err := server.Serve(l); err != nil {
		log.Fatal(err)
	}
}
