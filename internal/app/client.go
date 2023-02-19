package app

import (
	"google.golang.org/grpc"
	"imageclient/internal/config"
	"imageclient/internal/flags"
	"log"
)

func Run() {
	defer func(FileServer *grpc.ClientConn) {
		if err := FileServer.Close(); err != nil {
			log.Fatal(err)
		}
	}(config.FileServer)
	defer func(ListServer *grpc.ClientConn) {
		if err := ListServer.Close(); err != nil {
			log.Fatal(err)
		}
	}(config.ListServer)
	flags.Parse()
}
