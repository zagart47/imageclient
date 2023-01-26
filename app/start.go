package app

import (
	"imageclient/config"
	"imageclient/parseflags"
)

func Start() {

	defer config.ConnFile.Close()
	defer config.ConnList.Close()
	parseflags.ParseFlags()

}
