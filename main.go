package main

import (
	"errors"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/k0kubun/pp"
	"github.com/ktr0731/evans-demo/api"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type User struct {
	id                  string
	firstName, lastName string
	gender              api.Gender
}

type UserService struct {
	store map[string]*User
}

func (s *UserService) RegisterUsers(ctx context.Context, req *api.RegisterUsersRequest) (*api.RegisterUsersResponse, error) {
	for _, user := range req.GetUsers() {
		id := uuid.NewV4().String()
		s.store[id] = &User{
			id:        id,
			firstName: user.GetFirstName(),
			lastName:  user.GetLastName(),
			gender:    user.GetGender(),
		}
	}
	return &api.RegisterUsersResponse{
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
	pp.Println(req.GetId())
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
	api.RegisterUserServiceServer(server, &UserService{store: map[string]*User{}})
	if err := server.Serve(l); err != nil {
		log.Fatal(err)
	}
}
