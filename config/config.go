package config

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ConnFile, _ = grpc.Dial("0.0.0.0:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
var ConnList, _ = grpc.Dial("0.0.0.0:81", grpc.WithTransportCredentials(insecure.NewCredentials()))
