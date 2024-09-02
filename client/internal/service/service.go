package service

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

type WalletStruct struct {
	Address   string
	Signature string
}

var Addresses = GenerateAddresses(10000)

func CreateAndSignWallet() (*WalletStruct, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	signature, err := SignData(privateKey, []byte(address))
	if err != nil {
		return nil, err
	}

	return &WalletStruct{
		Address:   address,
		Signature: hex.EncodeToString(signature),
	}, nil
}

func SignData(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return crypto.Sign(hash.Bytes(), privateKey)
}

func GenerateAddresses(n int) []string {
	addresses := make([]string, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Failed to generate private key: %v", err)
		}

		address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		addresses[i] = address
	}
	return addresses
}

func ChunkAddresses(addresses []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(addresses); i += chunkSize {
		end := i + chunkSize
		if end > len(addresses) {
			end = len(addresses)
		}
		chunks = append(chunks, addresses[i:end])
	}
	return chunks
}

func LoadConfig() string {
	viper.SetConfigName("config")
	viper.AddConfigPath("server/config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}
	return viper.GetString("server_port")
}
