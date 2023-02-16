package download

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"imageclient/config"
	"imageclient/file"
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
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	md := metadata.Pairs("filename", name)
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	fd := pb.NewFileServiceClient(config.ConnFile)
	downloadStream, err := fd.Download(mdCtx, &pb.DownloadRequest{})
	if err != nil {
		return nil, err
	}
	f := file.NewFile(name)
	for {
		req, err := downloadStream.Recv()
		if err == io.EOF {
			if err := os.WriteFile(f.Path, f.Buffer.Bytes(), 0644); err != nil {
				log.Println(err.Error())
			}
			break
		}
		if err != nil {
			log.Println(err.Error())
			break
		}
		f.Buffer.Write(req.GetFragment())
	}
	log.Printf("file downloaded: %s", f.Name)
	return nil, nil
}

func (lsclient FileListClient) GetFileList() {
	table, err := lsclient.lsclient.GetFiles(context.Background(), &pb.GetFilesRequest{})
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(table.GetInfo())
}
