package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisCli, err := GetRedisCli()
	if err != nil {
		fmt.Printf("连接Redis失败了，err:%s", err)
		return
	}
	//fmt.Println(redisCli)

	//
	//redisCLIURL, err := GetRedisCliByURL()
	//if err != nil {
	//	fmt.Printf("连接Redis失败了，err:%s", err)
	//	return
	//}
	//fmt.Println(redisCLIURL)

	//redisCLISSH, err := GetRedisCliSSH()
	//if err != nil {
	//	fmt.Printf("连接Redis失败了，err:%s", err)
	//	return
	//}
	//fmt.Println(redisCLISSH)

	// 获取值
	get := redisCli.Get(context.TODO(), "key222")
	fmt.Println(get.Val())
	fmt.Println(get.Err())

	//可以使用 Do() 方法执行尚不支持或者任意命令:
	val, err := redisCli.Do(context.TODO(), "get", "key").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		panic(err)
	}
	fmt.Println(val.(string))

	val, err = redisCli.Get(context.TODO(), "key2").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key不存在")
	case err != nil:
		fmt.Println("错误", err)
	case val == "":
		fmt.Println("值是空字符串")
	}

}
