package redis

import (
	"time"
	"wikipay-admin/tools/config"

	"github.com/go-redis/redis"
)

var (
	clusterClient *redis.ClusterClient
)

//集群
func Connect() {
	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:       config.RedisClusterConfig.Addrs,
		Password:    config.RedisClusterConfig.Password,
		DialTimeout: time.Second * time.Duration(config.RedisClusterConfig.Dialtimeout),
		PoolSize:    config.RedisClusterConfig.Poolsize,
	})

	// fmt.Println(config.RedisClusterConfig.Addrs)
	_, err := clusterClient.Ping().Result()
	if err != nil {
		panic("redis connect error")
	}
}

//
func ClusterClient() *redis.ClusterClient {
	return clusterClient
}
