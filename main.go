package main

import (
	"flag"
	"fmt"
	"log"


	"google.golang.org/grpc"
	
	"github.com/golang/glog"
)

var (
	grpcPort       = flag.String("grpc", "9000", "gRPC server port")

)

type envoyAuthzServerV3 struct{}

type envoyAuthzServer struct {
	grpcServer *grpc.Server
	grpcV3     *envoyAuthzServerV3
}



