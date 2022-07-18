package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// create user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// listen user channel
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}

// user online
func (u *User) Online() {
	// user online, add to OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// broad user online message
	u.server.BroadCast(u, "上線了 !!!")
}

// user offline
func (u *User) Offline() {
	// user online, remove to OnlineMap
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// broad user online message
	u.server.BroadCast(u, "下線了 888")
}

// user send message
func (u *User) DoMessage(msg string) {
	if msg == "who" { // list online users
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[Server] %s 在線上 ...", user.Name)
			u.sendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if strings.Contains(msg, "rename") && msg[:6] == "rename" { // rename:xxx
		newName := msg[7:]

		// check name is exist
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.sendMsg(fmt.Sprintf("The %s has been used", newName))
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.sendMsg("Rename success")
		}
	} else {
		u.server.BroadCast(u, msg)
	}
}

// send message to now user
func (u *User) sendMsg(msg string) {
	u.conn.Write([]byte(msg + "\n"))
}
