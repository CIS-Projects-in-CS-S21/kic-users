package main

import (
	"context"
	"flag"
	"fmt"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"github.com/gogo/googleapis/google/rpc"
	"log"
	"strings"
	"google.golang.org/genproto/googleapis/rpc/status"
)

const (
	checkHeader   = "x-ext-authz"
	allowedValue  = "allow"
	resultHeader  = "x-ext-authz-check-result"
	resultAllowed = "allowed"
	resultDenied  = "denied"
)

var (
	serviceAccount = flag.String("allow_service_account", "a", "allowed service account, matched against the service account in the source principal from the client certificate")
	httpPort       = flag.String("http", "8000", "HTTP server port")
	grpcPort       = flag.String("grpc", "9000", "gRPC server port")
	denyBody       = fmt.Sprintf("denied by ext_authz for not found header `%s: %s` in the request", checkHeader, allowedValue)
)

type extAuthzServerV3 struct{}

type ExtAuthzServer struct {
	grpcServer *grpc.Server
	grpcV3     *extAuthzServerV3
	// For test only
	grpcPort chan int
}

type AuthService struct {
	logger *zap.SugaredLogger
}

type User struct {
	username string
	password string
}

type DataBase interface {
	CheckUser(user *User) (bool, error)
}

// Check implements gRPC v3 check request.
func (s *extAuthzServerV3) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	l := fmt.Sprintf("%s%s, attributes: %v\n",
		request.GetAttributes().GetRequest().GetHttp().GetHost(),
		request.GetAttributes().GetRequest().GetHttp().GetPath(),
		request.GetAttributes())
	if allowedValue == request.GetAttributes().GetRequest().GetHttp().GetHeaders()[checkHeader]  || strings.HasSuffix(request.GetAttributes().Source.Principal, "/sa/" + *serviceAccount) {
		log.Printf("[gRPCv3][allowed]: %s", l)
		return &authv3.CheckResponse{
			HttpResponse: &authv3.CheckResponse_OkResponse{
				OkResponse: &authv3.OkHttpResponse{
					Headers: []*corev3.HeaderValueOption{
						{
						},
					},
				},
			},
			Status: &status.Status{Code: int32(rpc.OK)},
		}, nil
	}

	log.Printf("[gRPCv3][denied]: %s", l)
	return &authv3.CheckResponse{
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden},
				Body:   denyBody,
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   resultHeader,
							Value: resultDenied,
						},
					},
				},
			},
		},
		Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)},
	}, nil
}

func(user *User) CheckUser() (bool, error) {
	return true,nil
}



