package test

import (
	"encoding/gob"
	"testing"
	"time"

	"gorpc.demo.com/client"
	"gorpc.demo.com/server"
	"gorpc.demo.com/service"
)

func init() {
	gob.Register(map[string]float64{})
}

func startTestServer() {
	go func() {
		s := server.NewServer(":1234")
		s.Register("Calculator", &service.Calculator{})
		s.Serve()
	}()
	// 等待 server 启动监听
	time.Sleep(200 * time.Millisecond)
}

func TestCalculatorAdd(t *testing.T) {
	startTestServer()
	cli, err := client.Dial(":1234")
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	params := map[string]float64{"A": 1, "B": 2}
	var result float64
	err = cli.Call("Calculator", "Add", params, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != 3 {
		t.Fatalf("expected 3, got %v", result)
	}
}

func startHTTPTestServer() {
	go func() {
		s := server.NewServer(":1234")
		s.Register("Calculator", &service.Calculator{})
		err := s.StartHTTP(":8080")
		if err != nil {
			return
		}
	}()
	time.Sleep(200 * time.Millisecond)
}

func TestCalculatorAddHTTP(t *testing.T) {
	//startHTTPTestServer()
	hcli := client.NewHTTPClient("http://localhost:8080/rpc")
	params := map[string]float64{"A": 2, "B": 5}
	var result float64
	err := hcli.CallHTTP("Calculator.Add", params, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != 7 {
		t.Fatalf("expected 7, got %v", result)
	}
}

func TestCalculatorAddHTTPAsync(t *testing.T) {
	//startHTTPTestServer()
	hcli := client.NewHTTPClient("http://localhost:8080/rpc")
	params := map[string]float64{"A": 3, "B": 4}
	var result float64
	errCh := hcli.AsyncCallHTTP("Calculator.Add", params, &result)
	err := <-errCh
	if err != nil {
		t.Fatal(err)
	}
	if result != 7 {
		t.Fatalf("expected 7, got %v", result)
	}
}
