package server

import (
	"fmt"
	"testing"
)


func TestUsersService_GenerateJWT(t *testing.T) {
	s := NewUsersService(nil, nil)
	tok, err := s.GenerateJWT(20)

	if err != nil {
		t.Errorf("error generating: %v", err)
	}

	res, err := s.DecodeJWT(tok)

	if err != nil {
		t.Errorf("error generating: %v", err)
	}

	fmt.Printf("%v %v\n", res, err)

	id, _ := res.Get("uid")

	if string(id) != 20 {
		t.Errorf("Key did not parse to the proper value got: %v", id)
	}


}