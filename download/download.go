package download

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"imageclient/config"
	pb "imageclient/pkg/proto"
	"io"
	"io/ioutil"
	"log"
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

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: pb.NewFileServiceClient(conn),
	}
}

func (c Client) Download(name string) (pb.FileService_DownloadClient, error) {
	l := pb.NewFileServiceClient(config.Conn)
	fileStreamResponse, err := l.Download(context.TODO(), &pb.DownloadRequest{Filename: name})
	if err != nil {
		return nil, err
	}

	for {
		req, err := fileStreamResponse.Recv()
		if err == io.EOF {
			log.Println("received all fragments")
			break
		}
		if err != nil {
			log.Println("error receiving fragments")
			break
		}
		ioutil.WriteFile("files/"+name, req.GetFragment(), 0644)
	}
	return nil, nil
}

func (c Client) GetFileList() FileList {
	res, err := c.client.GetFiles(context.Background(), &pb.GetFilesRequest{})
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
