package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/kic/users/pkg/database"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/status"
)

const (
	authHeader = "Authorization"
	denyBody = "Bad credentials"
)

func (s *UsersService) DecodeJWT(payload string) (jwt.Token, error) {
	token, err := jwt.Parse(
		[]byte(payload),
		// Tell the parser that you want to use this keyset
		jwt.WithKeySet(s.keyset),
		jwt.UseDefaultKey(true),
	)

	return token, err
}

func (s *UsersService) GenerateJWT(userID int64) (string, error) {
	t := jwt.New()
	err := t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	if err != nil {
		return "", err
	}

	err = t.Set("uid", strconv.FormatInt(userID, 10))
	if err != nil {
		return "", err
	}

	key, _ := s.keyset.Get(0)

	signed, err := jwt.Sign(t, jwa.HS256, key)

	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func (s *UsersService) ValidateUser(username, password string) (bool, error) {
	res, err := s.db.GetUser(context.TODO(), &database.UserModel{
		Username: username,
	})

	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))

	if err != nil {
		return false, nil
	}
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
		s.logger.Infof("[gRPCv3][allowed]: %s", l)
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

	s.logger.Infof("[gRPCv3][denied]: %s", l)
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