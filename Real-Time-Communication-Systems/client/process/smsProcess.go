package process

import (
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
)

type SmsProcess struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMsg(content string) (err error) {
	//1. 创建一个Msg
	var msg message.Message
	msg.Type = message.SmsMsgType

	//2. 创建一个SmsMsg
	var smsMsg message.SmsMsg
	smsMsg.Content = content
	smsMsg.UserId = CurUser.UserId
	smsMsg.UserStatus = CurUser.UserStatus

	//3.序列化smsMsg
	data, err := json.Marshal(smsMsg)
	if err != nil {
		fmt.Println("序列化失败在SendGroup err=", err.Error())
		return
	}
	msg.Data = string(data)

	//4.序列化smsMsg
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化失败在SendGroup err=", err.Error())
		return
	}
	//5. 将msg发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMsg err=", err)
		return
	}
	return
}
