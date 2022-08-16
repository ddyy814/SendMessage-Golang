package message

import (
	_"fmt"
)

// 确定一些消息类型
const (
	LoginMsgType = "LoginMsg"
	LoginResMsgType = "LoginResMsg"
	RegisterMsgType = "RegiesterMsg"
	RegisterResMsgType = "RegisterResMsg"
	NotifyUserStatusMsgType = "NotifyUserStatusMsg"
)


// 定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline 
	UserBusyStatus
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
	UserId []int			//增加字段， 保存用户ID的切片
	Error string `json:"error"` //返回错误信息
}


type RegisterMsg struct {
	User User `json:"user"` //类型就是User结构体
}

type RegisterResMsg struct {
	Code int `json:"code"` //返回状态码 400表示该用户已经占有， 200表示注册成功
	Error string `json:"error"` //返回错误信息
}

// 推送上线用户状态变化消息
type NotifyUserStatusMsg struct {
	UserId int `json: "userId"`
	Status int `json: "status"`
}