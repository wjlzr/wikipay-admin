package utils

import (
	"github.com/go-redis/redis"
	"io/ioutil"
)

const CountryKey = "countrys_key"

type Flags struct {
	CountryCode string `json:"countryCode"`
	TwoCharCode string `json:"twoCharCode"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Flag        string `json:"flag"`
	FlagC       string `json:"flag-c"`
}

//导入国家旗帜到redis中
func SetFlags(client *redis.ClusterClient, path string) error {
	var (
		err error
		fs  []byte
	)

	if client.Get(CountryKey).Val() == "" {
		fs, err = ioutil.ReadFile(path)
		if err == nil {
			_, err = client.Set(CountryKey, fs, 0).Result()
		}
	}
	return err
}

//获取国家信息
func FindFlag(countryCode string) *Flags {
	//val := myredis.ClusterClient().Get(CountryKey).Val()
	//
	//if val != "" {
	//	var flags []Flags
	//	if err := json.Unmarshal([]byte(val), &flags); err == nil {
	//		if len(flags) > 0 {
	//			for _, f := range flags {
	//				if f.CountryCode == countryCode {
	//					return &f
	//				}
	//			}
	//		}
	//	}
	//}
	//var flag Flags
	//flag = flag{Name: "aa", Code: "bb"}
	var s1 *Flags = &Flags{Name: "aa", Code: "bb"}
	return s1
	//return nil
}
