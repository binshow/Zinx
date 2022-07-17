package ziface

// -------------------------------------------
// @file          : iserver.go
// @author        : binshow
// @time          : 2022/7/17 11:51 AM
// @description   :	定义服务器接口
// -------------------------------------------

type IServer interface {
	Start()
	Stop()
	Server()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(router IRouter)
}

