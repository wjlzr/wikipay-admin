package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var cfgDatabase *viper.Viper
var cfgApplication *viper.Viper
var cfgJwt *viper.Viper
var cfgLog *viper.Viper
var cfgRedisCluster *viper.Viper

//func init() {
//	InitConfig("settings.dev")
//}

//载入配置文件
func ConfigSetup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("config not found settings.database")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgJwt = viper.Sub("settings.jwt")
	if cfgJwt == nil {
		panic("config not found settings.jwt")
	}
	JwtConfig = InitJwt(cfgJwt)

	cfgLog = viper.Sub("settings.log")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	LogConfig = InitLog(cfgLog)

	cfgRedisCluster = viper.Sub("settings.rediscluster")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	RedisClusterConfig = InitRedisCluster(cfgRedisCluster)
}

func SetApplicationIsInit() {
	SetConfig("./config", "settings.application.isInit", false)
}

func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	viper.WriteConfig()
}
