package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"test1/tools"
)

func main() {
	http.HandleFunc("/submitOrder", submitOrderHandler)

	http.HandleFunc("/submitLogin", submitLoginHandler)

	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Println("服务器正在 :8080 端口运行...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func submitLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "不允许的方法", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取请求体时出错", http.StatusInternalServerError)
		return
	}

	var loginData tools.LoginData
	err = json.Unmarshal(body, &loginData)

	Username, err := strconv.Atoi(loginData.Username)
	if err != nil {
		errorMsg := fmt.Sprintf("无法将字符串转换为整数: %v\n请求体内容: %s", err, string(body))
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	loginData.Username = strconv.Itoa(Username)

	// 处理订单
	fmt.Printf("登录信息：%+v\n", loginData)

	// 将订单信息保存到 Redis
	err = tools.SaveOrderToRedis1(loginData)
	if err != nil {
		log.Printf("保存订单到 Redis 时出错: %v", err)
		http.Error(w, "保存订单到 Redis 时出错", http.StatusInternalServerError)
		return
	}

	// 响应客户端
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "订单已成功接收"})
}

func submitOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "不允许的方法", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取请求体时出错", http.StatusInternalServerError)
		return
	}

	var order tools.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		errorMsg := fmt.Sprintf("解码 JSON 时出错: %v\n请求体内容: %s", err, string(body))
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(order.Quantity)
	if err != nil {
		errorMsg := fmt.Sprintf("无法将字符串转换为整数: %v\n请求体内容: %s", err, string(body))
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	order.Quantity = strconv.Itoa(quantity)

	// 处理订单
	fmt.Printf("收到订单：%+v\n", order)

	// 将订单信息保存到 Redis
	err = tools.SaveOrderToRedis(order)
	if err != nil {
		log.Printf("保存订单到 Redis 时出错: %v", err)
		http.Error(w, "保存订单到 Redis 时出错", http.StatusInternalServerError)
		return
	}

	// 响应客户端
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "订单已成功接收"})
}
