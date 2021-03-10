package server

import (
	"context"
	"fmt"
	"log"


	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/genproto/googleapis/rpc/status"

)

const (
	authHeader = "Authorization"
	denyBody = "Bad credentials"
)

func (s *UsersService) ValidateUser(username, password string) (bool, error) {
	return true, nil
}

func parseCredentialsFromHeader(header string) (string, string) {
	return "", ""
}

// Check implements gRPC v3 check request.
func (s *UsersService) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	l := fmt.Sprintf("%s%s, attributes: %v\n",
		request.GetAttributes().GetRequest().GetHttp().GetHost(),
		request.GetAttributes().GetRequest().GetHttp().GetPath(),
		request.GetAttributes())


	creds := request.GetAttributes().GetRequest().GetHttp().GetHeaders()[authHeader]

	username, password := parseCredentialsFromHeader(creds)

	approve, err := s.ValidateUser(username, password)

	if err == nil && approve {
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
					},
				},
			},
		},
		Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)},
	}, nil
}