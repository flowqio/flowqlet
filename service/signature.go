package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/stevensu1977/toolbox/crypto"
)

func Signature(code, nonce, ts string) bool {
	waitSign := []string{apiToken, nonce, ts}
	sort.Strings(waitSign)
	sign := crypto.Sha1([]byte(strings.Join(waitSign, "")))
	fmt.Println(sign, nonce, ts)
	if code != sign {
		return false
	}
	return true
}
