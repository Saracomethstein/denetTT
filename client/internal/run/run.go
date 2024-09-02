package run

import (
	"context"
	"denetTT/account"
	"denetTT/client/internal/service"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

var tokens = []string{
	"0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599", // WBTC
	"0x628F76eAB0C1298F7a24d337bBbF1ef8A1Ea6A24", // XRP
	"0xB8c77482e45F1F44dE1745F52C74426C631bDD52", // BNB
	"0x582d872A1B094FC48F5DE31D3B73F2D9bE47def1", // TONCOIN
	"0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE", // SHIB
}

func GetAccount(client account.AccountServiceClient, wallet *service.WalletStruct) {
	start := time.Now()

	req := &account.GetAccountRequest{
		EthereumAddress: wallet.Address,
		CryptoSignature: wallet.Signature,
	}

	resp, err := client.GetAccount(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}

	fmt.Printf("GetAccount Response: Gastoken Balance: %s, Wallet Nonce: %s\n", resp.GastokenBalance, resp.WalletNonce)
	fmt.Printf("GetAccount took %s\n", time.Since(start))
}

func GetAccounts(client account.AccountServiceClient, totalRequests int) {
	fmt.Printf("\nRunning test with %d total requests...\n", totalRequests)
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	defer cancel()

	stream, err := client.GetAccounts(ctx)
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	addressChunks := service.ChunkAddresses(service.Addresses[:totalRequests], totalRequests/len(tokens))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, token := range tokens {
		for _, chunk := range addressChunks {
			wg.Add(1)
			sem <- struct{}{}

			go func(chunk []string, token string) {
				defer wg.Done()
				defer func() { <-sem }()

				req := &account.GetAccountsRequest{
					EthereumAddresses: chunk,
					Erc20TokenAddress: token,
				}

				if err := stream.Send(req); err != nil {
					log.Printf("Failed to send request for token %s: %v", token, err)
				}
			}(chunk, token)
		}
	}

	go func() {
		wg.Wait()
		stream.CloseSend()
	}()

	handleResponses(stream, totalRequests)

	fmt.Printf("Test with %d total requests took %s\n", totalRequests, time.Since(start))
}

func handleResponses(stream account.AccountService_GetAccountsClient, totalRequests int) {
	var resultWg sync.WaitGroup
	resultChan := make(chan *account.GetAccountsResponse, totalRequests)

	resultWg.Add(1)
	go func() {
		defer resultWg.Done()
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Failed to receive response: %v", err)
				break
			}
			resultChan <- resp
		}
		close(resultChan)
	}()

	go func() {
		for resp := range resultChan {
			fmt.Printf("Address: %s, Balance: %d\n", resp.EthereumAddress, resp.Erc20Balance)
		}
	}()
	resultWg.Wait()
}
