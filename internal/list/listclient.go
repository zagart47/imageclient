package list

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "imageclient/proto"
	"log"
)

type File struct {
	FileName string
	Created  string
	Updated  string
}

type FileList []File

type Client struct {
	client pb.ListServiceClient
}

func New(conn grpc.ClientConnInterface) Client {
	return Client{client: pb.NewListServiceClient(conn)}
}

func (c Client) GetFileList() {
	table, err := c.client.GetFiles(context.Background(), &pb.GetFilesRequest{})
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(table.GetInfo())
}
