package model

import (
	"testing"

	"github.com/stevensu1977/toolbox/crypto"
)

var passSHA1 = "e35bece6c5e6e0e86ca51d0440e92282a9d6ac8a"

func TestNewUser(t *testing.T) {
	user := NewUser()
	user.Password = crypto.Sha1("welcome1")
	if user.Password != passSHA1 {
		t.Fatalf("%s not match %s", user.Password, passSHA1)
	}
	t.Log("test correct")

}
