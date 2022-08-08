package main

import (
	"fmt"
	"net" // the package for socket
)

func process(conn net.Conn) {
	// get data from Client
	defer conn.Close() // close conn

	for {
		// create new slice everytime
		buf := make([]byte, 1024)
		// Waiting for Client send data, if no write from Client, goroutine will be blocked
	//	fmt.Printf("Server is waiting for message...... from Client=%s\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf) // read data from Conn
		if err != nil {
			fmt.Printf("Client exit err=%v!!", err)
			return 
		}

		// display data of Client on server endpoint
		fmt.Print(string(buf[:n]))
	}

}

func main() {
	fmt.Println("Server start to listen.....")
	listen, err := net.Listen("tcp", "127.0.0.1:8888") // tcp protocol and listening 8888 port

	if err != nil {
		fmt.Println("listen error=", err)
		return
	}

	defer listen.Close() // close listen port

	for {
		fmt.Println("Waiting for connecting...")
		// waiting for connect from client
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept() err=%v", err)
		} else {
			fmt.Printf("Accept() successful conn=%v clientIP=%v\n", conn, conn.RemoteAddr().String())
		}

		// create goroutine to serve client
		go process(conn)
	}

	// fmt.Printf("listen suc=%v\n", listen)
}
