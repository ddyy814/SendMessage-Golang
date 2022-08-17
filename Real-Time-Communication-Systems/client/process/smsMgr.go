package process

import (
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
)

func outputGroupMsg(msg *message.Message) {
	//显示即可
	//1. 反序列化
	var smsMsg message.SmsMsg
	err := json.Unmarshal([]byte(msg.Data), &smsMsg)
	if err != nil {
		fmt.Println("反序列化失败err=", err.Error())
		return
	}
	
	//显示信息
	info := fmt.Sprintf("用户ID:\t%d  对大家说:\t%s ", smsMsg.UserId, smsMsg.Content)
	fmt.Println(info)
	fmt.Println()
}