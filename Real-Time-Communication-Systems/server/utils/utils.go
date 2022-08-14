package utils

import (
	"fmt"
	"net"
	"encoding/json"
	"go-code/SendMessageProject/Real-Time-Communication-Systems/common/message"
	"encoding/binary"
)


//将这些方法关联到结构体里
type Transfer struct {
	Conn net.Conn
	Buf [8096]byte

}

func (this *Transfer) ReadPkg() (msg message.Message, err error) {
	// buf := make([]byte, 8096)
	fmt.Println("读取客服换发送的数据...")
	n, err := this.Conn.Read(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Read err=", err)
		return
	}
	//根据bug[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	//根据pkeLen读取消息内容
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err=", err)
		return
	}

	//把pkgLen反序列化 -> message.Message
	//一定要加&， 要不加msg是空的
	err = json.Unmarshal(this.Buf[:pkgLen], &msg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}

	//发送这个data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}