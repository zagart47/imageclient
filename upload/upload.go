package upload

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"imageclient/config"
	pb "imageclient/pkg/proto"
	"io"
	"log"
	"os"
	"time"
)

type Client struct {
	client pb.FileServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: pb.NewFileServiceClient(conn),
	}
}

func (c Client) Upload(ctx context.Context, file string) (string, error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1000*time.Second))
	defer cancel()

	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("cannot open file (%s)", err.Error())
	}
	fileStat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", fmt.Errorf("cannot upload context (%s)", err.Error())
	}

	fileName := fileStat.Name()

	md := metadata.Pairs("filename", fileName)
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	buffer := make([]byte, 1024)
	fu := pb.NewFileServiceClient(config.ConnFile)
	uploadStream, err := fu.Upload(mdCtx)

	for {
		n, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("buffer reading error (%s)", err.Error())
		}
		req := &pb.UploadRequest{Fragment: buffer[:n]}
		if err := uploadStream.Send(req); err != nil {
			log.Fatal(err.Error())
		}
	}
	_, err = uploadStream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("cannot send file to server (%s)", err.Error())
	}
	log.Println("file uploaded:", file)
	return fmt.Sprintf("file uploaded %s", file), nil
}
