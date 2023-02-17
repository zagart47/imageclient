package listclient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb2 "imageclient/proto"
	"log"
)

type File struct {
	FileName string
	Created  string
	Updated  string
}

type FileList []File

type Client struct {
	lc pb2.ListServiceClient
}

func New(conn grpc.ClientConnInterface) Client {
	return Client{lc: pb2.NewListServiceClient(conn)}
}

func (c Client) GetFileList() {
	table, err := c.lc.GetFiles(context.Background(), &pb2.GetFilesRequest{})
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(table.GetInfo())
}
