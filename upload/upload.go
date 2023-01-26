package upload

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "imageclient/pkg/proto"
	"io"
	"log"
	"os"
	"strings"
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
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("cannot open file (%s)", err.Error())
	}

	stream, err := c.client.Upload(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot upload context (%s)", err.Error())
	}
	TrimFileNamePrefix(&file)
	req := &pb.UploadRequest{
		Filename: file,
	}
	err = stream.Send(req)
	if err != nil {
		return "", fmt.Errorf("cannot send filename to stream (%s)", err.Error())
	}
	reader := bufio.NewReader(f)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("buffer reading error (%s)", err.Error())
		}
		req := &pb.UploadRequest{
			Fragment: buffer[:n],
		}
		err = stream.Send(req)
		if err != nil {
			return "", fmt.Errorf("cannot send filebody to stream (%s)", err.Error())
		}

	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return "", fmt.Errorf("cannot send file to server (%s)", err.Error())
	}
	log.Println("file uploaded:", file)
	return fmt.Sprintf("file uploaded %s", file), nil
}

func TrimFileNamePrefix(filename *string) {
	for {
		if strings.Contains(*filename, "/") || strings.Contains(*filename, "\\") {
			*filename = (*filename)[1:]
		} else {
			break
		}
	}
}
