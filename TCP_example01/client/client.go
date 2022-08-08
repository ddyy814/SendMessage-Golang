package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}

	// Send data to server, then Exit
	reader := bufio.NewReader(os.Stdin) // standard input

	for {
		// read data from user input, then send to server
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read string failed err=", err)
		}

		// if user input = Exit, break
		line = strings.Trim(line, " \r\n")
		if line == "exit" {
			fmt.Println("Client exit!!!")
			break
		}

		// Send line to server
		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Conn.Write err=", err)
		}
	}
}
