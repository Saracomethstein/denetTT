package handlers

import (
	"context"
	"denetTT/account"
	"denetTT/server/internal/server/service"
	"fmt"
	"io"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Server struct {
	account.UnimplementedAccountServiceServer
	Client *ethclient.Client
}

func (s *Server) GetAccount(ctx context.Context, req *account.GetAccountRequest) (*account.GetAccountResponse, error) {
	addres := common.HexToAddress(req.EthereumAddress)

	if !service.VerifySignature(addres, req.CryptoSignature) {
		return nil, fmt.Errorf("invalid signature")
	}

	balance, nonce, err := service.GetBalanceAndNonce(s.Client, addres)
	if err != nil {
		return nil, err
	}

	return &account.GetAccountResponse{
		GastokenBalance: strconv.FormatUint(balance, 10),
		WalletNonce:     strconv.FormatUint(nonce, 10),
	}, nil
}
func (s *Server) GetAccounts(stream account.AccountService_GetAccountsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		for _, ethAddress := range req.EthereumAddresses {
			address := common.HexToAddress(ethAddress)
			tokenAddress := common.HexToAddress(req.Erc20TokenAddress)

			balance, err := service.GetERC20Balance(s.Client, address, tokenAddress)
			if err != nil {
				return err
			}

			if err := stream.Send(&account.GetAccountsResponse{
				EthereumAddress: ethAddress,
				Erc20Balance:    balance,
			}); err != nil {
				return err
			}
		}
	}
}
