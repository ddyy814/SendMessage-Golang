package main

import (
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/process"
	"net"
	"io"
)

// 创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

//编写一个serverProcessMsg函数
// 功能：根据客户端发送消息种类不同， 决定调用哪个函数处理
func (this *Processor) serverProcessMsg(msg *message.Message) (err error) {

	switch msg.Type {
	case message.LoginMsgType:
		//处理登陆的逻辑
		//创建一个user Process实例
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessLogin(msg)
	case message.RegisterMsgType:
		//处理注册
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessRegister(msg)
	case message.SmsMsgType:
		//创建一个SmsProcess,
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMsg(msg)
	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

func (this *Processor) process2 ()  (err error){
	//循环读客户端发送的信息
	for {
		//我们将读取数据包直接封装成一个函数readPkg()，返回Message,Err
		tf := &utils.Transfer{
			Conn : this.Conn,
		}
		msg, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也正常退出...")
				return err
			} else {
				fmt.Println("readPkg fail err=", err)
				return err
			}
		}
		// fmt.Println("msg=", msg)
		err = this.serverProcessMsg( &msg)
		if err != nil {
			return err
		}
	}
}