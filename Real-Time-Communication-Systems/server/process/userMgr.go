package process2

import (
	"fmt"
)


// 因为UserMgr，在服务器端只有一个
// 因为在很多地方用， 所以我们定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对UserMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//删除
func (this *UserMgr) DeleteOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

//查询
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

//根据ID返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//从map中取出一个值，带检测方式
	up, ok := this.onlineUsers[userId]
	if !ok { // 说明要查找的这个用户，当前不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}