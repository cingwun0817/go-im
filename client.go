package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1. Public Mode")
	fmt.Println("2. Private Mode")
	fmt.Println("3. Rename")
	fmt.Println("0. Quit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag

		return true
	} else {
		fmt.Println(">>>> Please input 0-3 menu numbers")
		return false
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {

		}

		switch c.flag {
		case 1:
			fmt.Println("Choose Public Mode")

			break
		case 2:
			fmt.Println("Choose Private Mode")

			break
		case 3:
			// fmt.Println("Choose Rename Mode")

			c.Rename()

			break
		}
	}
}

func (c *Client) Rename() bool {
	fmt.Println(">>>>> Input new name")

	fmt.Scanln(&c.Name)

	sendMsg := "rename:" + c.Name + "\n"

	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.conn)
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "setting server ip (default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "setting server port (default: 8888)")
}

func main() {
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> client connection failed ...")
		return
	}

	fmt.Printf(">>>>> %s %d, client connection success ...\n", serverIp, serverPort)

	// response server message to stdout
	go client.DealResponse()

	client.Run()

	// // block forever
	// select {}
}
