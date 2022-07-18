package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// create server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

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

	// listen message channel
	go s.ListenMessage()

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
	fmt.Printf("Connect success: %s %s\n", conn.RemoteAddr().String(), conn.RemoteAddr().Network())

	user := NewUser(conn)

	// user online, add to OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// broad user online message
	s.BroadCast(user, "上線了!!!")

	// receive user message
	go func() {
		buffer := make([]byte, 1024)

		for {
			n, err := conn.Read(buffer)

			if n == 0 {
				s.BroadCast(user, "下線了QQQ")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
				return
			}

			// get message (filter \n)
			msg := string(buffer[0 : n-1])

			// broad cast msg
			s.BroadCast(user, msg)
		}
	}()

	// block forever
	select {}
}

// broad cast msg
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s] %s: %s", user.Addr, user.Name, msg)

	s.Message <- sendMsg
}

// listen message
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		// to online users
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}
