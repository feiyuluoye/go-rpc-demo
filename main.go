package main

import (
	"encoding/gob"
	"fmt"

	"gorpc.demo.com/server"
	"gorpc.demo.com/service"
)

func main() {
	gob.Register(map[string]float64{})

	s := server.NewServer(":1234")
	s.Register("Calculator", &service.Calculator{})

	// 启动 TCP-RPC 服务
	go func() {
		fmt.Println("TCP-RPC 服务已启动，监听 :1234")
		s.Serve()
	}()

	// 启动 HTTP-RPC 服务
	go func() {
		fmt.Println("HTTP-RPC 服务已启动，监听 :8080/rpc")
		err := s.StartHTTP(":8080")
		if err != nil {
			fmt.Println("HTTP-RPC 启动失败:", err)
		}
	}()

	select {} // 阻塞主线程
}
