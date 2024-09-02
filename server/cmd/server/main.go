package main

import (
	"denetTT/account"
	"denetTT/server/internal/server/handlers"
	"denetTT/server/internal/server/service"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	infuraURL := viper.GetString("infura_url")
	serverPort := viper.GetString("server_port")

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
