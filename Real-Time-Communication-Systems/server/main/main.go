package main

import (
	"fmt"
	"net"
)

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里也需要延迟关闭
	defer conn.Close()

	//调用processor
	processor := &Processor{
		Conn :conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯协程错误", err)
		return 
	}
}

func main() {
	//提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("Listen出错 err=", err)
		return
	}
	//一旦监听成功，等待客户端连接服务器
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept出错")
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}
