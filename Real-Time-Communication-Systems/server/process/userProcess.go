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
