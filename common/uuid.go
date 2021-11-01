package common

import (
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

//获取商户密钥
func GenerateKey() string {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	return strings.Replace(uuid.Must(uuid.NewV4(), nil).String(), "-", "", -1)
}
