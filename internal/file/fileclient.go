package file

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"imageclient/internal/config"
	"imageclient/internal/model"
	pb "imageclient/proto"
	"io"
	"log"
	"os"
	"time"
)

type Client struct {
	client pb.FileServiceClient
}

func New(conn grpc.ClientConnInterface) Client {
	return Client{client: pb.NewFileServiceClient(conn)}
}

func (c Client) Download(name string) (pb.FileService_DownloadClient, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	md := metadata.Pairs("filename", name)
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	downloadClient := New(config.FileServer)
	stream, err := downloadClient.client.Download(mdCtx, &pb.DownloadRequest{})
	if err != nil {
		return nil, err
	}
	f := model.NewFile(name)
	for {
		req, err := stream.Recv()
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
		return "", fmt.Errorf("cannot open file (%w)", err)
	}
	stat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	fileName := stat.Name()

	md := metadata.Pairs("filename", fileName)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	buffer := make([]byte, 1024)
	fileUpload := pb.NewFileServiceClient(config.FileServer)
	stream, err := fileUpload.Upload(mdCtx)
	if err != nil {
		return "", err
	}

	for {
		bytes, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("buffer reading error (%w)", err)
		}
		req := &pb.UploadRequest{Fragment: buffer[:bytes]}
		if err := stream.Send(req); err != nil {
			log.Fatal(err.Error())
		}
	}
	if _, err = stream.CloseAndRecv(); err != nil {
		return "", fmt.Errorf("cannot send file to server (%w)", err)
	}

	fmt.Println("file uploaded:", fileName)
	return fmt.Sprintf("file uploaded %s", fileName), nil
}
