package app

import (
	"imageclient/config"
	"imageclient/parseflags"
)

func Start() {

	defer config.Conn.Close()
	parseflags.ParseFlags()

}
