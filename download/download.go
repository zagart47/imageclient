package download

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"imageclient/config"
	"imageclient/model"
	pb "imageclient/pkg/proto"
	"io"
	"log"
	"os"
	"time"
)

type File struct {
	FileName string
	Created  string
	Updated  string
}

type FileList []File

type Client struct {
	client pb.FileServiceClient
}

type FileListClient struct {
	lsclient pb.ListServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: pb.NewFileServiceClient(conn),
	}
}

func NewListServiceClient(conn grpc.ClientConnInterface) FileListClient {
	return FileListClient{lsclient: pb.NewListServiceClient(conn)}
}

func (c Client) Download(name string) (pb.FileService_DownloadClient, error) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	md := metadata.Pairs("filename", name)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	fd := pb.NewFileServiceClient(config.ConnFile)
	downloadStream, err := fd.Download(mdCtx, &pb.DownloadRequest{})
	if err != nil {
		return nil, err
	}
	f := model.NewFile(name)
	for {
		req, err := downloadStream.Recv()
		if err == io.EOF {
			err1 := os.WriteFile(f.Path, f.Buffer.Bytes(), 0644)
			if err1 != nil {
				log.Println(err.Error())
			}
			break
		}
		f.Buffer.Write(req.GetFragment())
	}
	return nil, err
}

func (lsclient FileListClient) GetFileList() {
	table, err := lsclient.lsclient.GetFiles(context.Background(), &pb.GetFilesRequest{})
	if err != nil {
	}
	fmt.Println(table.GetInfo())
}
