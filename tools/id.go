package tools

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	"wikipay-admin/redis"
)

const (
	secretKey = "3Emd#(09)wer#@900EPpoKLK!2($"

	imUnuse = "{im_easemobs}_unuse"
	imUsed  = "{im_easemobs}_used"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

//获取用户id号
func GenerateUserId() int64 {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	rx := RandInt64(1000000000, 9999999999)
	return rx
}

//随机数-int64
func RandInt64(min, max int64) int64 {
	return min + r.Int63n(max-min)
}

//加密
func Encrypt(pwd string) string {
	if pwd == "" {
		return ""
	}

	saltedPwd := fmt.Sprintf("%s_%s", pwd, secretKey)
	encrypedPwd := fmt.Sprintf("%x", sha256.Sum256([]byte(saltedPwd)))
	return encrypedPwd
}

//获取
func GetImInfo() []string {
	count := redis.ClusterClient().SCard(imUnuse).Val()
	if count > 0 {
		im := redis.ClusterClient().SRandMember(imUnuse).Val()
		if im != "" {
			strs := strings.Split(im, "|")
			if len(strs) == 2 {
				err := redis.ClusterClient().SMove(imUnuse, imUsed, im).Err()
				if err != nil {
					return nil
				}
				return strs
			}
		}
	}
	return nil
}
