package config

import (
	"google.golang.org/grpc"
)

var Conn, err = grpc.Dial("0.0.0.0:80", grpc.WithInsecure())
