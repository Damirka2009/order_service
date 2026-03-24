package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPC_PORT string
	HTTP_PORT string
	NETWORK   string
}

func Load() *Config {
	godotenv.Load()
	gRPC_port := os.Getenv("GRPC_PORT")
	http_port := os.Getenv("HTTP_PORT")
	network := os.Getenv("NETWORK")
	return &Config{
		GRPC_PORT: gRPC_port,
		HTTP_PORT: http_port,
		NETWORK:   network,
	}
}
