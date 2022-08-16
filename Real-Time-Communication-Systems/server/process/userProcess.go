package process2

import (
	"encoding/json"
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/model"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn

	//增加一个字段，表示该Conn是哪个用户
	UserId int
}

//编写通知所有在线用
// userId要通知其他的在线用户，我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers,然后一个个发送
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知
		up.NotifyMeOnlineUser(userId)
	}
}

func (this *UserProcess) NotifyMeOnlineUser(userId int) {
	//组装NotifyUserStatusMsg
	var msg message.Message
	msg.Type = message.NotifyUserStatusMsgType

	var notifyUserStatusMsg message.NotifyUserStatusMsg
	notifyUserStatusMsg.UserId = userId
	notifyUserStatusMsg.Status = message.UserOnline

	//将notifyUserStatusMsg序列化
	data, err := json.Marshal(notifyUserStatusMsg)
	if err != nil {
		fmt.Println("序列化出错err=", err)
		return
	}
	//将序列化后的notifyUserStatusMsg赋值给msg.Data
	msg.Data = string(data)

	//对msg再次序列化， 准备发送
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化出错err=", err)
		return
	}

	// 发送， 创建transfer实例， 发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnlineUser err=",err)
		return
	}

}


func (this *UserProcess) ServerProcessRegister(msg *message.Message) (err error) {

	// 1.先从msg 中取出msg.Data, 并直接反序列化 RegisterMsg
	var registerMsg message.RegisterMsg
	err = json.Unmarshal([]byte(msg.Data), &registerMsg)
	if err != nil {
		fmt.Println("反序列化失败。", err)
		return
	}

	// 1. 先声明一个 response msg
	var resMsg message.Message
	resMsg.Type = message.RegisterResMsgType
	// 2. 再声明一个RegisterResMsg
	var registeResrMsg message.RegisterResMsg

	// 我们需要到redis数据库去完成注册
	// 1. 使用model.MyUserDao去redis验证
	err = model.MyUserDao.Register(&registerMsg.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registeResrMsg.Code = 505
			registeResrMsg.Error = model.ERROR_USER_EXISTS.Error()
		}else {
			registeResrMsg.Code = 506
			registeResrMsg.Error = "注册发生未知错误..."
		}
	}else {
		registeResrMsg.Code = 200
	}

	data, err := json.Marshal(registeResrMsg)
	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}

	// 4.将data赋值给 resMsg
	resMsg.Data = string(data)

	// 5. 对resMsg进行序列化， 准备发送
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}
	// 6. 发送data 我们将其封装到writePkg
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// 编写一个函数serverProcessLogin函数， 专门处理登陆请求
func (this *UserProcess) ServerProcessLogin(msg *message.Message) (err error) {
	//核心代码
	// 1.先从msg 中取出msg.Data, 并直接反序列化LoginMsg
	var loginMsg message.LoginMsg
	err = json.Unmarshal([]byte(msg.Data), &loginMsg)
	if err != nil {
		fmt.Println("反序列化失败。", err)
		return
	}

	// 1. 先声明一个 response msg
	var resMsg message.Message
	resMsg.Type = message.LoginResMsgType

	// 2. 再声明一个LoginReMsg
	var loginResMsg message.LoginResMsg

	// 我们需要到redis数据库去完成验证
	// 1. 使用model.MyUserDao去redis验证
	user, err := model.MyUserDao.Login(loginMsg.UserId, loginMsg.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMsg.Code = 500
			loginResMsg.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMsg.Code = 403
			loginResMsg.Error = err.Error()
		} else {
			loginResMsg.Code = 505
			loginResMsg.Error = "服务器内部错误..."
		}
	} else {
		loginResMsg.Code = 200
		//这里用户登陆成功， 我们把该登陆成功的用户放入到UserMgr中
		//将登陆成功的用户ID，赋给this
		this.UserId = loginMsg.UserId
		//通知其他在线用户，我上线了
		userMgr.AddOnlineUser(this)
		this.NotifyOthersOnlineUser(loginMsg.UserId)
		//将当前在线用户的ID，放入到loginResMsg.UserId
		//遍历userMgr.onlineUsers
		for id,_ := range userMgr.onlineUsers {
			loginResMsg.UserId = append(loginResMsg.UserId, id)
		}
		fmt.Println(user, "loggedin")
	}

	// // 如果用户的ID为100， 密码是12345，就是合法的
	// if loginMsg.UserId == 100 && loginMsg.UserPwd == "123456" {
	// 	//合法
	// 	loginResMsg.Code = 200
	// } else {
	// 	//不合法
	// 	loginResMsg.Code = 500 //表示用户不存在
	// 	loginResMsg.Error = "该用户不存在，请注册在使用"
	// }

	// 3. 将loginResMsg序列化
	data, err := json.Marshal(loginResMsg)

	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}
	// 4.将data赋值给 resMsg
	resMsg.Data = string(data)

	// 5. 对resMsg进行序列化， 准备发送
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal failed", err)
		return
	}
	// 6. 发送data 我们将其封装到writePkg
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
