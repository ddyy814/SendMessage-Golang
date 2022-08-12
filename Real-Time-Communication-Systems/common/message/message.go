package message

import (
	_"fmt"
)

// 确定一些消息类型
const (
	LoginMsgType = "LoginMsg"
	LoginResMsgType = "LoginResMsg"
	RegisterMsgType = "RegiesterMsg"
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"`// 消息内容
}


//定义两个消息
type LoginMsg struct {
	UserId int `json:"userId"`	//用户ID
	UserPwd string `json:"userPwd"` //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMsg struct {
	Code int `json:"code"` //返回状态码 500表示用户未注册， 200表示登陆成功
	Error string `json:"error"` //返回错误信息
}


type RegisterMsg struct {

}
