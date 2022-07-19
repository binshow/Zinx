package ziface

// -------------------------------------------
// @file          : imsghandler.go
// @author        : binshow
// @time          : 2022/7/19 9:06 AM
// @description   : 消息管理抽象层
// -------------------------------------------


/*
	消息管理抽象层
*/
type IMsgHandle interface{
	DoMsgHandler(request IRequest)			//马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter)	//为消息添加具体的处理逻辑
}