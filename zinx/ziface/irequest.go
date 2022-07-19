package ziface

// -------------------------------------------
// @file          : irequest.go
// @author        : binshow
// @time          : 2022/7/17 3:21 PM
// @description   : 将客户端的所有请求信息封装成 Request
// -------------------------------------------

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface{
	GetConnection() IConnection	//获取请求连接信息
	GetData() []byte			//获取请求消息的数据
	GetMsgID() uint32
}