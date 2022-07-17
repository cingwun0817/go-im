package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// create server
func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port}

	return server
}

// start server
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}

}

// handler
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("Connect success")
}
