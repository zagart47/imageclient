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
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1000*time.Second))
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

func (lsclient FileListClient) GetFileList() FileList {
	res, err := lsclient.lsclient.GetFiles(context.Background(), &pb.GetFilesRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fl := FileList{}

	for _, v := range res.Info {
		fl = append(fl, File{
			FileName: v.FileName,
			Created:  v.Created,
			Updated:  v.Updated,
		})
	}
	Decorate()
	fmt.Println("|     Имя файла      |    Дата создания    |   Дата обновления   |")
	Decorate()
	for _, v := range fl {

		fmt.Printf("|%s| %s | %s |\n", Fitting(v.FileName, 20), v.Created, Fitting(v.Updated, 19))
	}
	Decorate()
	return fl
}

func Decorate() {
	fmt.Println("+--------------------+---------------------+---------------------+")
}

func Fitting(s string, n int) string {
	for len(s) < n {
		s = s + " "
		if len(s) == n {
			break
		}
		s = " " + s
		if len(s) == n {
			break
		}
	}
	return s
}
