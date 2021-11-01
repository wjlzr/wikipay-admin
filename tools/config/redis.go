package config

import "github.com/spf13/viper"

var RedisClusterConfig = new(RedisCluster)

type RedisCluster struct {
	Addrs       []string
	Password    string
	Dialtimeout int64
	Poolsize    int
}

func InitRedisCluster(cfg *viper.Viper) *RedisCluster {
	return &RedisCluster{
		Addrs:       cfg.GetStringSlice("addrs"),
		Password:    cfg.GetString("password"),
		Dialtimeout: cfg.GetInt64("dialtimeout"),
		Poolsize:    cfg.GetInt("poolsize"),
	}
}
