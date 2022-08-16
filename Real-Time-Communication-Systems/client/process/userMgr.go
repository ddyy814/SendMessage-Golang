package process

import (
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
)

// 客户端要维护的mp
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

//在客户端显示当前在线用户
func outputOnlineUser() {
	//遍历 onlineUsers
	fmt.Println("当前在线用户列表:")
	for id, _ := range onlineUsers {
		fmt.Println("用户ID:\t", id)
	}
}

//编写一个方法， 处理返回的 NotifyUserStatusMsg
func updateUserStatus(notifyUserStatusMsg *message.NotifyUserStatusMsg) {

	//适当优化
	user, ok := onlineUsers[notifyUserStatusMsg.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMsg.UserId,
		}
	}

	user.UserStatus = notifyUserStatusMsg.Status
	
	onlineUsers[notifyUserStatusMsg.UserId] = user
	outputOnlineUser()
}
