package process

import (
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
	"net"
	"os"
)

//显示登陆成功后的界面
func ShowMenu() {
	fmt.Println("--------------------登陆成功--------------------")
	fmt.Println("--------------------1. 显示在线用户列表--------------------")
	fmt.Println("--------------------2. 发送消息--------------------")
	fmt.Println("--------------------3. 信息列表--------------------")
	fmt.Println("--------------------4. 退出系统--------------------")
	fmt.Println("请选择(1-4):")

	var key int
	fmt.Scanf("%d\n", &key)

	switch key {
	case 1:
		// fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入的选项不正确...")
	}
}


func serverProcessMsg(conn net.Conn) {
	//创建一个transfer，不停读取服务器发送消息
	tf := &utils.Transfer{
		Conn : conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		msg, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("Readpkg err", err)
			return
		}
		//如果读取到消息，下一步处理逻辑
		// fmt.Printf("msg=%v\n", msg)
		switch msg.Type {
		case message.NotifyUserStatusMsgType : //有人上线了
			// 1. 取出NotifyUserStatusMsg
			var notifyUserStatusMsg message.NotifyUserStatusMsg
			json.Unmarshal([]byte(msg.Data), &notifyUserStatusMsg)

			// 2. 把这个用户的信息状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMsg)
		default:
			fmt.Println("u服务器端返回了未知的消息类型")
		}

	}
}