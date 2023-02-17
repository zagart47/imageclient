package fileclient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"imageclient/config"
	"imageclient/file"
	pb2 "imageclient/proto"
	"io"
	"log"
	"os"
	"time"
)

type Client struct {
	client pb2.FileServiceClient
}

func New(conn grpc.ClientConnInterface) Client {
	return Client{client: pb2.NewFileServiceClient(conn)}
}

func (c Client) Download(name string) (pb2.FileService_DownloadClient, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	md := metadata.Pairs("filename", name)
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	downloadClient := New(config.FileServer)
	downloadStream, err := downloadClient.client.Download(mdCtx, &pb2.DownloadRequest{})
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

func (c Client) Upload(file string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()

	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("cannot open file (%s)", err.Error())
	}
	fileStat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	fileName := fileStat.Name()

	md := metadata.Pairs("filename", fileName)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	buffer := make([]byte, 1024)
	fu := pb2.NewFileServiceClient(config.FileServer)
	uploadStream, err := fu.Upload(mdCtx)

	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("buffer reading error (%s)", err.Error())
		}
		req := &pb2.UploadRequest{Fragment: buffer[:n]}
		if err := uploadStream.Send(req); err != nil {
			log.Fatal(err.Error())
		}
	}
	if _, err = uploadStream.CloseAndRecv(); err != nil {
		return "", fmt.Errorf("cannot send file to server (%s)", err.Error())
	}

	fmt.Println("file uploaded:", fileName)
	return fmt.Sprintf("file uploaded %s", fileName), nil
}
