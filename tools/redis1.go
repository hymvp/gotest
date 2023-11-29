// redis1.go

package tools

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Order struct {
	ProductName string `json:"productName"`
	Quantity    string `json:"quantity"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var redisClient *redis.Client

func init() {
	// 初始化 Redis 客户端
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 你的 Redis 服务器地址
		Password: "",               // 如果有密码，填写密码
		DB:       0,                // 使用默认的数据库
	})

	// 检查是否成功连接到 Redis
	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("无法连接到 Redis: %v", err))
	}
	fmt.Println("成功连接到 Redis:", pong)
}

// SaveOrderToRedis 将订单信息保存到 Redis
func SaveOrderToRedis(order Order) error {
	// 将订单信息转换为 JSON 格式
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// 使用 Redis 客户端将订单信息保存到 Redis 中
	err = redisClient.Set("order:"+order.ProductName, orderJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func SaveOrderToRedis1(loginData LoginData) error {
	// 将订单信息转换为 JSON 格式
	orderJSON, err := json.Marshal(loginData)
	if err != nil {
		return err
	}

	// 使用 Redis 客户端将订单信息保存到 Redis 中
	err = redisClient.Set("order:"+loginData.Username, orderJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
