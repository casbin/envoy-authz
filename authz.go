package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	//"encoding/json"
	"strings"
	"encoding/base64"

	//"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"	
	//"github.com/golang/glog"
	"github.com/casbin/casbin"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/gogo/googleapis/google/rpc"
)

var (
	grpcport       = flag.String("grpc", "9000", "grpcport")	

)


type AuthorizationServer struct{}

/*func (a *AuthorizationServer) GetUserName(req *http.Request) string {
	username, _, _ := req.BasicAuth()
	return username
}
*/

func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	log.Println(">>> Server Performing authorization check!")
	
	e:= casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
/*	if err != nil {
		glog.Errorf("Filed to load the policies: %v", err)
		
	}
*/
/*	user, _, _ := req.BasicAuth()
	method := req.Method
	path := req.URL.Path
*/
//	user :=  a.GetUserName(req) 
//	user, _, _ := json.Unmarshal(req.Attributes.Request.Http.Headers["authorization"].BasicAuth())

	method := req.Attributes.Request.Http.Method
	path := req.Attributes.Request.Http.Path
	authHeader, ok := req.Attributes.Request.Http.Headers["authorization"]
	if !ok {
		log.Fatalf("failed to receive headers")
	}
	var splitToken []string
	splitToken = strings.Split(authHeader, "Bearer ")
	if len(splitToken) == 2 {
		token := splitToken[1]
		tokenbyte, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			panic("malformed input")
		}
		tokenStr :=  string(tokenbyte[:])
		tokenvalue := strings.Split(tokenStr, ",")
		username := tokenvalue[1]
		if e.Enforce(username, path, method) {
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
		}else {
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
	}
	return &auth.CheckResponse{
		Status: &rpcstatus.Status{
			Code: int32(rpc.UNAUTHENTICATED),
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: "Authorization Header malformed or not provided",
			},
		},
	}, nil
	
	//tokenvalue := strings.Split(tokenStr, ",")
	//username := tokenvalue[1]
	/*
	if e.Enforce(username, path, method) {
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
	}else {
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

	}*/
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
