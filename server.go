package server

import (
	"fmt"
)

func server() {
	// 启动服务器
	fmt.Println("Server is starting on port 8080...")
	user := new(struct{})
	// TODO: 增加鉴权逻辑
	if user == nil {
		fmt.Println("Forbidden!")
		return
	}
}
