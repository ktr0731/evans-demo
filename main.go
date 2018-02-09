package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/net/context"

	"github.com/ktr0731/evans-demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterService struct{}

func (s *GreeterService) SayHello(ctx context.Context, req *api.HelloRequest) (*api.HelloResponse, error) {
	var msg string
	switch req.GetLanguage() {
	case api.Language_ENGLISH:
		msg = "Hello, %s also %s!"
	case api.Language_JAPANESE:
		msg = "こんにちは, %s と %s!"
	default:
		return nil, status.Error(codes.InvalidArgument, "unknown language type")
	}
	you := makeName(req.GetYou())
	theirs := make([]string, 0, len(req.GetTheirs()))
	for _, p := range req.GetTheirs() {
		theirs = append(theirs, makeName(p))
	}
	return &api.HelloResponse{
		Message: fmt.Sprintf(msg, you, strings.Join(theirs, ", ")),
	}, nil
}

func makeName(p *api.Person) string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

func main() {
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	api.RegisterGreeterServer(server, &GreeterService{})
	if err := server.Serve(l); err != nil {
		log.Fatal(err)
	}
}
