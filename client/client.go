package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"time"

	"gorpc.demo.com/protocol"
)

type Client struct {
	conn net.Conn
}

func Dial(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Call(service, method string, params interface{}, result interface{}) error {
	req := &protocol.Request{
		Service: service,
		Method:  method,
		Params:  params,
	}
	err := protocol.WriteRequest(c.conn, req)
	if err != nil {
		return err
	}
	resp, err := protocol.ReadResponse(c.conn)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return fmt.Errorf(resp.Error)
	}
	// 假设 result 是指针
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(resp.Result))
	return nil
}

// CallJSON 通过 TCP 发送 JSON-RPC 2.0 请求
func (c *Client) CallJSON(method string, params interface{}, result interface{}) error {
	jsonReq := &protocol.JSONRequest{
		JSONRPC: "2.0",
		Method:  method, // 格式: Service.Method
		Params:  params,
		ID:      1,
	}
	err := protocol.WriteJSONRequest(c.conn, jsonReq)
	if err != nil {
		return err
	}
	jsonResp, err := protocol.ReadJSONResponse(c.conn)
	if err != nil {
		return err
	}
	if jsonResp.Error != nil && jsonResp.Error != "" {
		return fmt.Errorf("rpc error: %v", jsonResp.Error)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(jsonResp.Result))
	return nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

// HTTPClient 支持 HTTP-RPC (JSON-RPC 2.0)
type HTTPClient struct {
	url    string
	client *http.Client
}

func NewHTTPClient(url string) *HTTPClient {
	return &HTTPClient{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// CallHTTP 同步调用
func (c *HTTPClient) CallHTTP(method string, params interface{}, result interface{}) error {
	jsonReq := &protocol.JSONRequest{
		JSONRPC: "2.0",
		Method:  method, // 格式: Service.Method
		Params:  params,
		ID:      1,
	}
	buf := new(bytes.Buffer)
	err := protocol.WriteJSONRequest(buf, jsonReq)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.url, "application/json", buf)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	jsonResp, err := protocol.ReadJSONResponse(resp.Body)
	if err != nil {
		return err
	}
	if jsonResp.Error != nil && jsonResp.Error != "" {
		return fmt.Errorf("rpc error: %v", jsonResp.Error)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(jsonResp.Result))
	return nil
}

// AsyncCallHTTP 异步调用，返回结果 channel
func (c *HTTPClient) AsyncCallHTTP(method string, params interface{}, result interface{}) <-chan error {
	ch := make(chan error, 1)
	go func() {
		err := c.CallHTTP(method, params, result)
		ch <- err
	}()
	return ch
}
