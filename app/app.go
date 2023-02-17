package app

import (
	"imageclient/config"
	"imageclient/parseflags"
)

func Start() {

	defer config.FileServer.Close()
	defer config.ListServer.Close()
	parseflags.ParseFlags()

}
