package flags

import (
	"errors"
	"flag"
	"imageclient/internal/config"
	"imageclient/internal/file"
	"imageclient/internal/list"
	"log"
)

type Connect struct {
	upload   file.Client
	download file.Client
	list     list.Client
}

func NewConnect() *Connect {
	return &Connect{
		upload:   file.New(config.FileServer),
		download: file.New(config.FileServer),
		list:     list.New(config.ListServer),
	}
}

func Parse() {
	if err := Perform(); err != nil {
		log.Fatal(err)
	}
}

type Flag struct {
	download string
	upload   string
	list     bool
}

func NewFlags() Flag {
	return Flag{
		download: "",
		upload:   "",
		list:     false,
	}
}

var flags = NewFlags()

func Perform() error {
	connect := NewConnect()
	flag.StringVar(&flags.download, "dl", "", "download file")
	flag.StringVar(&flags.upload, "ul", "", "upload file")
	flag.BoolVar(&flags.list, "ls", false, "list files")
	flag.Parse()

	switch {
	case flags.list:
		connect.list.GetFileList()
		return nil
	case len(flags.download) > 0:
		if _, err := connect.download.Download(flags.download); err != nil {
			return err
		}
	case len(flags.upload) > 0:
		if _, err := connect.upload.Upload(flags.upload); err != nil {
			return err
		}
	default:
		return errors.New("operation not allowed")
	}
	return nil

}
