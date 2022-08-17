package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
	"net"
	"os"
)

type UserProcess struct {
}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延迟关闭
	defer conn.Close()

	//2. 准备通过conn发送消息给服务器
	var msg message.Message
	msg.Type = message.RegisterMsgType

	//3. 创建一个LoginMsg结构体
	var registerMsg message.RegisterMsg
	registerMsg.User.UserId = userId
	registerMsg.User.UserPwd = userPwd
	registerMsg.User.UserName = userName

	//4. registerMsgMsg序列化, data是切片，所以要勇敢string转换一下
	data, err := json.Marshal(registerMsg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5. data赋给msg.Data
	msg.Data = string(data)

	//6. 将msg序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
	}

	msg, err = tf.ReadPkg() // msg = regiestResMsg

	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将msg的data反序列化成RegisterResMsg
	var registerResMsg message.RegisterResMsg
	err = json.Unmarshal([]byte(msg.Data), &registerResMsg)
	if registerResMsg.Code == 200 {
		fmt.Println("注册成功，重新登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMsg.Error)
		os.Exit(0)
	}
	return
}

//写一个函数完成登陆校验
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

	// 下一步开始定协议..
	// fmt.Printf("userId = %d userPwd= %s", userId, userPwd)

	// return nil

	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延迟关闭
	defer conn.Close()

	//2. 准备通过conn发送消息给服务器
	var msg message.Message
	msg.Type = message.LoginMsgType

	//3. 创建一个LoginMsg结构体
	var loginMsg message.LoginMsg
	loginMsg.UserId = userId
	loginMsg.UserPwd = userPwd

	//4. loginMsg序列化, data是切片，所以要勇敢string转换一下
	data, err := json.Marshal(loginMsg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5. data赋给msg.Data
	msg.Data = string(data)

	//6. 将msg序列化
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7. data就是我们要发送的消息
	//7.1 先把data长度发送给服务器
	//先获取data长度，转成一个表示长度的byte的切片
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
	// fmt.Printf("客户端发送数据长度=%d 内容=%s 成功... ", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	// time.Sleep(20 * time.Second)
	// fmt.Println("Sleep 20..")
	//这里还需要处理服务器端返回的消息

	//创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	msg, err = tf.ReadPkg()

	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//将msg的data反序列化成LoginResMsg
	var loginResMsg message.LoginResMsg
	err = json.Unmarshal([]byte(msg.Data), &loginResMsg)
	if loginResMsg.Code == 200 {
		// fmt.Println("用户登陆成功")
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		//可以显示当前在线用户列表，遍历loginResMsg.UsersId
		fmt.Println("当前在线用户列表如下:")
		for _, v := range loginResMsg.UserId {

			//如果我们要求不显示自己在线，下面我们增加一个代码
			if v == userId {
				continue
			}
			fmt.Println("用户ID:\t", v)
			//完成客户端 onlineUsers的初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		// 这里我们需要在客户端启动一个协程
		//该协程保持和服务器端的通讯，如果服务器油数据推送给客户端
		//则接收并显示在客户端终端
		go serverProcessMsg(conn)

		// 1. 显示登陆成功后菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMsg.Error)
	}
	return
}
