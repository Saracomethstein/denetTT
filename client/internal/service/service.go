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

var Addresses = GenerateWallets(10000)

func GenerateWallets(n int) []string {
	addresses := make([]string, n)
	for i := 0; i < n; i++ {
		address, err := generateAddress()
		if err != nil {
			log.Printf("Failed to generate address: %v", err)
			continue
		}
		addresses[i] = address
	}
	return addresses
}

func GenerateAndSignWallet() (*WalletStruct, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	address, err := deriveAddress(privateKey)
	if err != nil {
		return nil, err
	}

	signature, err := signData(privateKey, []byte(address))
	if err != nil {
		return nil, err
	}

	return &WalletStruct{
		Address:   address,
		Signature: hex.EncodeToString(signature),
	}, nil
}

func generateAddress() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	return deriveAddress(privateKey)
}

func deriveAddress(privateKey *ecdsa.PrivateKey) (string, error) {
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKey).Hex(), nil
}

func signData(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return crypto.Sign(hash.Bytes(), privateKey)
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
		log.Fatalf("Error reading config file: %v", err)
	}
	return viper.GetString("server_port")
}
