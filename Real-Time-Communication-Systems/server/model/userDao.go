package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
)

// 我们在服务器启动后，就初始化一个userDao的实例
// 把它做成一个全局变量，在需要和redis操作时，直接使用即可
var (
	MyUserDao *UserDao
)

//定义一个UserDao结构体
//完成对User结构体操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式，创建一个userDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao){
	userDao = &UserDao{
		pool: pool,
	}
	return 
}

// 1. 根据用户ID， 返回一个user的实例+error
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//通过给定的ID去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		//error
		if err == redis.ErrNil { //表示在users 哈希中， 没有找到对应ID
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}
	// 需要把res反序列化成User实例
	json.Unmarshal([]byte(res), user)

	if err != nil {
		fmt.Println("反序列化出错err=", err)
		return
	}
	return
}

//完成登陆校验
// 1. Login 完成对用户的验证
// 2. 如果用户的ID和password都正确， 返回一个user实例
// 3. 如果用户的ID和password有错误，返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从UserDao链接池取出一个链接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	//证明这个用户获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return 
}

func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao链接池取出一个链接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	
	//这时， 说明ID在redis里还没有，所以可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return 
	}

	// 入库
	_, err = conn.Do("HSet", "users",user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return 
	}
	return
}