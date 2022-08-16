package message


// 定义一个用户的结构体
type User struct {
	//为了能够序列化和反序列化成功，我们必须要保证
	// 用户信息的json字符串的key和结构体的字段对应的tag名字一致
	UserId int `json: "userId"`
	UserPwd string `json: "userPwd"`
	UserName string `json: "userName"`
	UserStatus int `json:"userStatus"` //用户状态
}