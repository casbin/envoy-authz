package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"	
	"github.com/golang/glog"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/gogo/googleapis/google/rpc"
)

var (
	grpcport       = flag.String("grpc", "9000", "grpcport")	

)


type AuthorizationServer struct{}

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Server Performing authorization check!")
	
	e, err := casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
	if err != nil {
		glog.Errorf("Filed to load the policies: %v", err)
		return
	}

	user, _, _ := req.BasicAuth()
	method := req.Method
	path := req.URL.Path

	if e.Enforce(user, path, method) {
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(rpc.OK),
			},
			HttpResponse: &auth.CheckResponse_OkResponse{
				OkResponse: &auth.OkHttpResponse{
					Headers: []*core.HeaderValueOption{
						{
							Header: &core.HeaderValue{
								Key:   "casbin-authorized",
								Value: "allowed",
							},
						},
					},
				},
			},
		}, nil
	}
	else {
		return &auth.CheckResponse{
			Status: &rpcstatus.Status{
				Code: int32(rpc.PERMISSION_DENIED),
			},
			HttpResponse: &auth.CheckResponse_DeniedResponse{
				DeniedResponse: &auth.DeniedHttpResponse{
					Status: &envoy_type.HttpStatus{
						Code: envoy_type.StatusCode_Unauthorized,
					},
					Body: "PERMISSION_DENIED",
				},
			},
		}, nil

	}

func main(){
	flag.Parse()
	
	if *grpcport == "" {
		fmt.Fprintln(os.Stderr, "missing -grpcport flag ")
		flag.Usage()
		os.Exit(2)
	}
	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	auth.RegisterAuthorizationServer(s, &AuthorizationServer{})

	log.Printf("Starting gRPC Server at %s", *grpcport)
	s.Serve(lis)
}
