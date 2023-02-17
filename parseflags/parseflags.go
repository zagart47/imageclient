package parseflags

import (
	"errors"
	"flag"
	"imageclient/config"
	"imageclient/fileclient"
	"imageclient/listclient"
	"log"
)

type Connects struct {
	uploadClient   fileclient.Client
	downloadClient fileclient.Client
	listClient     listclient.Client
}

func NewConnect() *Connects {
	return &Connects{
		uploadClient:   fileclient.New(config.FileServer),
		downloadClient: fileclient.New(config.FileServer),
		listClient:     listclient.New(config.ListServer),
	}
}

func ParseFlags() {
	if err := Perform(); err != nil {
		log.Fatal(err)
	}
}

func Perform() error {
	c := NewConnect()
	downloadFlag := flag.String("dl", "", "download file")
	uploadFlag := flag.String("ul", "", "upload file")
	listFlag := flag.Bool("ls", false, "list files")
	flag.Parse()

	switch {
	case *listFlag:
		c.listClient.GetFileList()
		return nil
	case len(*downloadFlag) > 0:
		if _, err := c.downloadClient.Download(*downloadFlag); err != nil {
			return err
		}
	case len(*uploadFlag) > 0:
		if _, err := c.uploadClient.Upload(*uploadFlag); err != nil {
			return err
		}
	default:
		return errors.New("operation not allowed")
	}
	return nil

}
