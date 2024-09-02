package main

import (
	"denetTT/account"
	"denetTT/client/internal/run"
	"denetTT/client/internal/service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverAddress := service.LoadConfig()

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := account.NewAccountServiceClient(conn)

	wallet, err := service.CreateAndSignWallet()
	if err != nil {
		log.Fatalf("failed to create and sign wallet: %v", err)
	}

	run.GetAccount(client, wallet)
	run.GetAccounts(client, 100)
	run.GetAccounts(client, 1000)
	run.GetAccounts(client, 10000)
}
