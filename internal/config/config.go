package config

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var (
	_             = godotenv.Load()
	fileHost      = os.Getenv("fileHost")
	listHost      = os.Getenv("listHost")
	target        = NewTarget(fileHost, listHost)
	FileServer, _ = grpc.Dial(target.file, grpc.WithTransportCredentials(insecure.NewCredentials()))
	ListServer, _ = grpc.Dial(target.list, grpc.WithTransportCredentials(insecure.NewCredentials()))
)

type Target struct {
	file string
	list string
}

func NewTarget(file string, list string) Target {
	return Target{
		file: file,
		list: list,
	}
}
