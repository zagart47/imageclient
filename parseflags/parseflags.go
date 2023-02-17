package parseflags

import (
	"errors"
	"flag"
	"imageclient/config"
	"imageclient/download"
	"imageclient/upload"
	"log"
)

func ParseFlags() {
	err := Perform()
	if err != nil {
		log.Fatal(err)
	}
}

func Perform() error {
	up := upload.NewClient(config.ConnFile)
	down := download.NewClient(config.ConnFile)
	list := download.NewListServiceClient(config.ConnList)

	dl := flag.String("dl", "", "download file")
	ul := flag.String("ul", "", "upload file")
	ls := flag.Bool("ls", false, "list files")
	flag.Parse()

	switch {
	case *ls:
		list.GetFileList()
		return nil
	case len(*dl) > 0:
		if _, err := down.Download(*dl); err != nil {
			return err
		}
	case len(*ul) > 0:
		if _, err := up.Upload(*ul); err != nil {
			return err
		}
	default:
		return errors.New("operation not allowed")
	}
	return nil

}

/*operation := args["o"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	ls := args["ls"]

	if len(ls) >= 0 {
		list.GetFileList()
		return nil
	}

	fileName := args["dl"]
	if (operation == "download" || operation == "upload") && fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "u":
		item := args["f"]
		if item == "" {
			return errors.New("-file flag has to be specified")
		}
		if _, err := up.Upload(fileName); err != nil {
			return err
		}
		return nil

	case "ls":
		list.GetFileList()
		return nil

	case "download":
		id := args["f"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		if _, err := down.Download(fileName); err != nil {
			return err
		}
		return nil

	default:
		return errors.New("operation not allowed")
	}
}*/
