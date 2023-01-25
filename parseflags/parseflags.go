package parseflags

import (
	"context"
	"errors"
	"flag"
	"imageclient/config"
	"imageclient/download"
	"imageclient/upload"
)

func ParseFlags() {
	err := Perform(parseArgs())
	if err != nil {
		panic(err)
	}
}

type Arguments map[string]string

func parseArgs() Arguments {
	operationFlag := flag.String("o", "", "operation")
	fileNameFlag := flag.String("f", "", "fileName")
	flag.Parse()

	return Arguments{
		"o": *operationFlag,
		"f": *fileNameFlag,
	}
}

func Perform(args Arguments) error {
	up := upload.NewClient(config.Conn)
	down := download.NewClient(config.Conn)

	operation := args["o"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	fileName := args["f"]
	if (operation == "download" || operation == "upload") && fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "upload":
		item := args["f"]
		if item == "" {
			return errors.New("-file flag has to be specified")
		}
		_, err := up.Upload(context.Background(), fileName)
		if err != nil {
			return err
		}
		return nil

	case "list":
		down.GetFileList()
		return nil

	case "download":
		id := args["f"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		_, err := down.Download(fileName)
		if err != nil {
			return err
		}
		return nil

	default:
		return errors.New("Operation abcd not allowed!")
	}
}
