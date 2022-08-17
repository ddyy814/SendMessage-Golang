package model

import (
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"net"
)



// 因为在客户端， 我们很多地方用到curUser， 定义为全局的
type CurUser struct {
	Conn net.Conn
	message.User
}