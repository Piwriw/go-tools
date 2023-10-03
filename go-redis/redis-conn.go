package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// https://redis.uptrace.dev/zh/

func GetRedisCliSSH() (*redis.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            "sshusername",
		Auth:            []ssh.AuthMethod{ssh.Password("sshpassword")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", "127.0.0.1:22", sshConfig)
	if err != nil {
		return nil, err
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort("127.0.0.1", "6379"),
		Password: "redis123",
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		},
		// SSH不支持超时设置，在这里禁用
		ReadTimeout:  -2,
		WriteTimeout: -2,
	})

	err = redisCli.Set(context.TODO(), "key", "TestPing", 0).Err()
	if err != nil {
		return nil, err
	}

	val, err := redisCli.Get(context.TODO(), "key").Result()
	if err != nil {
		return nil, err
	}
	if val != "TestPing" {
		return nil, errors.New("好像存取Redis拿到一些奇怪的值")
	}
	return redisCli, nil

}

// GetRedisCliByURL URL形式连接redis
func GetRedisCliByURL() (*redis.Client, error) {
	//opt, err := redis.ParseURL("redis://:redis123@localhost:6379/<db>")
	url := "redis://:redis123@localhost:6379/"
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	redisCli := redis.NewClient(opts)
	err = redisCli.Set(context.TODO(), "key", "TestPing", 0).Err()
	if err != nil {
		return nil, err
	}

	val, err := redisCli.Get(context.TODO(), "key").Result()
	if err != nil {
		return nil, err
	}
	if val != "TestPing" {
		return nil, errors.New("好像存取Redis拿到一些奇怪的值")
	}
	return redisCli, nil
}

// GetRedisCli 通过配置项目获取Redis
func GetRedisCli() (*redis.Client, error) {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis123", // no password set
		DB:       0,          // use default DB
	})
	err := redisCli.Set(context.TODO(), "key", "TestPing", 0).Err()
	if err != nil {
		return nil, err
	}

	val, err := redisCli.Get(context.TODO(), "key").Result()
	if err != nil {
		return nil, err
	}
	if val != "TestPing" {
		return nil, errors.New("好像存取Redis拿到一些奇怪的值")
	}
	return redisCli, nil

}
