package redisutil

import (
	"gopkg.in/redis.v5"
)

var gClientPool = map[string]*redis.Client{}

// 新获取客户端
func NewClient(addr, password string) (pClient *redis.Client, err error) {
	if pClient = gClientPool[addr]; pClient == nil {
		pNewClient := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})

		if _, err = pNewClient.Ping().Result(); err == nil {
			pClient = pNewClient
			gClientPool[addr] = pNewClient
		}
	}

	return pClient, err
}

// 默认客户端
func DefaultClient() (client *redis.Client, err error) {
	addr := "localhost:6379" // default ip:port
	password := ""           // no password set

	return NewClient(addr, password)
}
