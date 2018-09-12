package service

import (
	"strconv"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {

	fixedNonce := "aabc"
	fixedTS := "1535347308472808000"
	fixedCode := "ad39eaf383229f1895d5a2a1459b39870bc16bde"

	isValid := Signature(fixedCode, fixedNonce, fixedTS)

	t.Logf("fixed check is : %v", isValid)

	isValid = Signature("123", "aabc", strconv.FormatInt(time.Now().UnixNano(), 10))
	if !isValid {
		t.Fatal("Signature failure ")
	}

	t.Log(isValid)
}
