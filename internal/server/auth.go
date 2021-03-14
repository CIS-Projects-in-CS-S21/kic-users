package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/kic/users/pkg/database"
)

const (
	authHeader = "authorization"
	denyBody = "Bad credentials"
	resultHeader  = "x-ext-authz-check-result"
	resultAllowed = "allowed"
)

func (s *UsersService) DecodeJWT(payload string) (jwt.Token, error) {
	token, err := jwt.Parse(
		[]byte(payload),
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
		s.logger.Debugf("Failed to get user from db to validate: %v", err)
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password))

	if err != nil {
		s.logger.Debugf("Failed to compare passwords: %v", err)
		return false, nil
	}
	s.logger.Debugf("User is valid, returning")
	return true, nil
}

func parseCredentialsFromHeader(header string) (string, error) {
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("invalid header format")
	}

	reqToken := strings.TrimSpace(splitToken[1])

	return reqToken, nil
}

// Check implements gRPC v3 check request.
func (s *UsersService) Check(ctx context.Context, request *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	l := fmt.Sprintf("%s%s, attributes: %v\n",
		request.GetAttributes().GetRequest().GetHttp().GetHost(),
		request.GetAttributes().GetRequest().GetHttp().GetPath(),
		request.GetAttributes())

	header := request.GetAttributes().GetRequest().GetHttp().GetHeaders()[authHeader]

	approve := true

	tok, err := parseCredentialsFromHeader(header)

	if err != nil {
		approve = false
	} else {
		_, err := s.DecodeJWT(tok)
		if err != nil {
			approve = false
		}
	}

	if approve {
		s.logger.Infof("[gRPCv3][allowed]: %s", l)
		return &authv3.CheckResponse{
			HttpResponse: &authv3.CheckResponse_OkResponse{
				OkResponse: &authv3.OkHttpResponse{
					Headers: []*corev3.HeaderValueOption{
						{
							Header: &corev3.HeaderValue{
								Key:   resultHeader,
								Value: resultAllowed,
							},
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