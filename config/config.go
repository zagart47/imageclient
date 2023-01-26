package config

import (
	"google.golang.org/grpc"
)

var ConnFile, _ = grpc.Dial("0.0.0.0:80", grpc.WithInsecure())
var ConnList, _ = grpc.Dial("0.0.0.0:81", grpc.WithInsecure())
