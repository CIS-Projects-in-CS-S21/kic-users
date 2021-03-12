package server

import (
	"os"
	"strconv"
	"testing"
)


func TestUsersService_GenerateJWTAndDecode(t *testing.T) {
	os.Setenv("SECRET_KEY", "I-hate-writing-tests")
	s := NewUsersService(nil, nil)
	tok, err := s.GenerateJWT(20)

	if err != nil {
		t.Errorf("error generating: %v", err)
	}

	res, err := s.DecodeJWT(tok)

	if err != nil {
		t.Errorf("error generating: %v", err)
	}

	id, _ := res.Get("uid")

	intID, err := strconv.Atoi(id.(string))

	if intID != 20 {
		t.Errorf("Key did not parse to the proper value got: %v", id)
	}
}