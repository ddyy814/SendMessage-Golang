package main

import (
	"fmt"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/client/process"
)

//定义两个变量，一个表示用户的ID，一个表示用户的密码
var userId int
var userPwd string
var userName string

func main() {

	// 接收用户选择
	var key int
	// 判断是否继续显示菜单
	// loop := true

	for true {
		fmt.Println("-----------------------欢迎登陆多人聊天系统---------------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的ID")
			fmt.Scanln(&userId)
			fmt.Println("请输入用户密码")
			fmt.Scanln(&userPwd)
			//创建一个UserProcess
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户ID:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(nickname):")
			fmt.Scanf("%s\n", &userName)
			// 2. 调用UserProcess， 完成注册请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			// loop = false
		default:
			fmt.Println("你的输入有误，请重新输入")
		}
	}
}
