package main

import (
	"denetTT/account"
	"denetTT/server/internal/server/handlers"
	"denetTT/server/internal/server/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	infuraURL, serverPort := service.LoadConfig()

	client := service.InitEthereumClient(infuraURL)

	list, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Error listening on port %s", serverPort)
	}

	s := grpc.NewServer()
	account.RegisterAccountServiceServer(s, &handlers.Server{Client: client})

	log.Printf("Starting server on port %s", serverPort)
	if err := s.Serve(list); err != nil {
		log.Fatalf("Error serving on port %s", serverPort)
	}
}
