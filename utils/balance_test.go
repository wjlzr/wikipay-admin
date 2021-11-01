package utils

import (
	"fmt"
	"testing"
)

func TestBalance(t *testing.T) {
	f := GetBtcBalance("1JdPfWF3NeZDtHAyBvT4SPpZtwNAzyyoGw")
	fmt.Println(f)
}
