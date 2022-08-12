package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"io"
	"net"
)

func readPkg(conn net.Conn) (msg message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客服换发送的数据...")
	n, err := conn.Read(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Read err=", err)
		return
	}
	//根据bug[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])
	//根据pkeLen读取消息内容
	n, err = conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err=", err)
		return
	}

	//把pkgLen反序列化 -> message.Message
	//一定要加&， 要不加msg是空的
	err = json.Unmarshal(buf[:pkgLen], &msg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	// 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	//发送这个data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}

// 编写一个函数serverProcessLogin函数， 专门处理登陆请求
func serverProcessLogin(conn net.Conn, msg *message.Message) (err error) {
	//核心代码
	// 1.先从msg 中取出msg.Data, 并直接反序列化LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Println("反序列化失败。", err)
		return
	}

	// 1. 先声明一个 response msg
	var resMsg message.Message
	resMsg.Type = message.LoginResMsgType

	// 2. 再声明一个LoginReMsg
	var loginResMsg message.LoginResMsg

	// 如果用户的ID为100， 密码是12345，就是合法的
	if loginMsg.UserId == 100 && loginMsg.UserPwd == "123456" {
		//合法
		loginResMsg.Code = 200
	} else {
		//不合法
		loginResMsg.Code = 500 //表示用户不存在
		loginResMsg.Error = "该用户不存在，请注册在使用"
	}

	// 3. 将loginResMsg序列化
	data, err := json.Marshal(loginResMsg)

	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}
	// 4.将data赋值给 resMsg
	resMsg.Data = string(data)

	// 5. 对resMsg进行序列化， 准备发送
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}
	// 6. 发送data 我们将其封装到writePkg
	err = writePkg(conn, data)
	return 
}

//编写一个serverProcessMsg函数
// 功能：根据客户端发送消息种类不同， 决定调用哪个函数处理
func serverProcessMsg(conn net.Conn, msg *message.Message) (err error) {

	switch msg.Type {
	case message.LoginMsgType:
		//处理登陆的逻辑
		err = serverProcessLogin(conn, msg)
	case message.RegisterMsgType:
		//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里也需要延迟关闭
	defer conn.Close()

	//循环读客户端发送的信息
	for {
		//我们将读取数据包直接封装成一个函数readPkg()，返回Message,Err

		msg, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也正常退出...")
				return
			} else {
				fmt.Println("readPkg fail err=", err)
				return
			}
		}
		// fmt.Println("msg=", msg)
		err = serverProcessMsg(conn, &msg)
		if err != nil {
			return
		}
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
