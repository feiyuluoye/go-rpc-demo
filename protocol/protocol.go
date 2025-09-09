package protocol

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

type Request struct {
	Service string
	Method  string
	Params  interface{}
}

type Response struct {
	Result interface{}
	Error  string
}

// JSON-RPC 2.0 结构体
// 参考 https://www.jsonrpc.org/specification

type JSONRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

type JSONResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// Gob 编解码
func WriteRequest(w io.Writer, req *Request) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(req)
}

func ReadRequest(r io.Reader) (*Request, error) {
	dec := gob.NewDecoder(r)
	req := &Request{}
	if err := dec.Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func WriteResponse(w io.Writer, resp *Response) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(resp)
}

func ReadResponse(r io.Reader) (*Response, error) {
	dec := gob.NewDecoder(r)
	resp := &Response{}
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// JSON-RPC 编解码
func WriteJSONRequest(w io.Writer, req *JSONRequest) error {
	enc := json.NewEncoder(w)
	return enc.Encode(req)
}

func ReadJSONRequest(r io.Reader) (*JSONRequest, error) {
	dec := json.NewDecoder(r)
	req := &JSONRequest{}
	if err := dec.Decode(req); err != nil {
		return nil, err
	}
	return req, nil
}

func WriteJSONResponse(w io.Writer, resp *JSONResponse) error {
	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}

func ReadJSONResponse(r io.Reader) (*JSONResponse, error) {
	dec := json.NewDecoder(r)
	resp := &JSONResponse{}
	if err := dec.Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
