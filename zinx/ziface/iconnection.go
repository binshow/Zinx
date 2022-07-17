package ziface

import "net"

// -------------------------------------------
// @file          : iconnection.go
// @author        : binshow
// @time          : 2022/7/17 12:21 PM
// @description   : 连接的抽象接口
// -------------------------------------------

type IConnection interface {

	Start()		// 启动链接

	Stop()

	GetTCPConnection()  *net.TCPConn	// 获取与当前连接绑定的 socket conn

	GetConnID() uint32

	RemoteAddr() net.Addr 	// 获取远程客户端的 TCP状态 ip + port

	//直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error


}

// 定义一个从处理连接业务的方法
type HandleFunc func(*net.TCPConn , []byte , int) error