package config

import (
	"google.golang.org/grpc"
)

var Conn, err = grpc.Dial("localhost:12223", grpc.WithInsecure())
