package process2

import (
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
	"net"
)

type SmsProcess struct {
}

//写方法转发消息
func (this *SmsProcess) SendGroupMsg(msg *message.Message) {
	//遍历服务器端的onlineUsers map[int]
	//将消息转发出去

	//取出msg内容
	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化出错 err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//这里需要过滤掉自己，不要发给自己
		if id == smsMsg.UserId {
			continue
		}
		this.SendMsgToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess) SendMsgToEachOnlineUser(data []byte, conn net.Conn) {
	//创建transfer，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
