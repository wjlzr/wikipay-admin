package tools

import (
	"fmt"
	"testing"
)

func TestId(t *testing.T) {
	pwd := Encrypt("wikifx123")
	fmt.Println(pwd)
}
