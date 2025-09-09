package server

import (
	"net"
	"net/http"
	"reflect"

	"gorpc.demo.com/protocol"
)

type Server struct {
	addr    string
	service map[string]interface{}
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		service: make(map[string]interface{}),
	}
}

func (s *Server) Register(name string, svc interface{}) {
	s.service[name] = svc
}

func (s *Server) Serve() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		req, err := protocol.ReadRequest(conn)
		if err != nil {
			return
		}
		svc, ok := s.service[req.Service]
		if !ok {
			protocol.WriteResponse(conn, &protocol.Response{Error: "service not found"})
			continue
		}
		method := reflect.ValueOf(svc).MethodByName(req.Method)
		if !method.IsValid() {
			protocol.WriteResponse(conn, &protocol.Response{Error: "method not found"})
			continue
		}
		in := []reflect.Value{reflect.ValueOf(req.Params)}
		out := method.Call(in)
		protocol.WriteResponse(conn, &protocol.Response{Result: out[0].Interface()})
	}
}

func (s *Server) StartHTTP(httpAddr string) error {
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		jsonReq, err := protocol.ReadJSONRequest(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := protocol.WriteJSONResponse(w, &protocol.JSONResponse{
				JSONRPC: "2.0",
				Error:   "invalid json request",
				ID:      nil,
			})
			if err != nil {
				return
			}
			return
		}
		// 解析 method: "Service.Method"
		serviceName, methodName := parseMethod(jsonReq.Method)
		svc, ok := s.service[serviceName]
		if !ok {
			err := protocol.WriteJSONResponse(w, &protocol.JSONResponse{
				JSONRPC: "2.0",
				Error:   "service not found",
				ID:      jsonReq.ID,
			})
			if err != nil {
				return
			}
			return
		}
		method := reflect.ValueOf(svc).MethodByName(methodName)
		if !method.IsValid() {
			err := protocol.WriteJSONResponse(w, &protocol.JSONResponse{
				JSONRPC: "2.0",
				Error:   "method not found",
				ID:      jsonReq.ID,
			})
			if err != nil {
				return
			}
			return
		}
		in := []reflect.Value{reflect.ValueOf(jsonReq.Params)}
		out := method.Call(in)
		err = protocol.WriteJSONResponse(w, &protocol.JSONResponse{
			JSONRPC: "2.0",
			Result:  out[0].Interface(),
			ID:      jsonReq.ID,
		})
		if err != nil {
			return
		}
	})
	return http.ListenAndServe(httpAddr, nil)
}

// parseMethod 解析 "Service.Method" 字符串
func parseMethod(full string) (service, method string) {
	for i := 0; i < len(full); i++ {
		if full[i] == '.' {
			return full[:i], full[i+1:]
		}
	}
	return full, ""
}
