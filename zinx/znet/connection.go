package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/zinx/ziface"
)

// -------------------------------------------
// @file          : connection.go
// @author        : binshow
// @time          : 2022/7/17 12:25 PM
// @description   :	封装connection
// -------------------------------------------

type Connection struct {
	// 当前链接的 socket tcp 套接字
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool  // 当前链接状态

	//该连接的处理方法router ， 等价于上面的 handleAPI
	Router  ziface.IRouter

	ExitBuffChan chan bool // 告知当前链接已经退出or停止

	msgChan  chan []byte // 1. 无缓冲管道，用于读、写两个goroutine之间的消息通信

}

// 连接的读业务方法
func (c *Connection) startReader()  {
	fmt.Println("Reader Goroutine is Running....")
	defer fmt.Println("connID = " , c.ConnID , "reader is exit ,remote add is " , c.RemoteAddr().String())
	defer c.Stop()

	for true {

		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg , err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn:c,
			//data:buf,
			msg: msg,
		}

		//从路由Routers 中找到注册绑定Conn的对应Handle
		go func (request ziface.IRequest) {
			//执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)


	}
}

// 2. 连接的写业务方法，将数据发送给客户端
func (c *Connection) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for  {

		select {
		case data := <-c.msgChan: // 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:	// 连接已经关闭了
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start() , connID = " , c.ConnID)
	// 启动从当前链接的读数据的业务
	// todo 启动从当前链接写数据的业务

	go c.startReader()
	go c.startWriter() // 4. 启动

	for {
		select {
		case <- c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}

}

func (c *Connection) Stop() {
	fmt.Println("Conn stop() , connID = " , c.ConnID)

	if c.isClosed == true {
		return
	}

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConnection()  *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return  errors.New("Pack error msg ")
	}

	//3. 写回客户端
	c.msgChan <- msg   //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return nil
}



func NewConnection(conn *net.TCPConn, connID uint32, isClosed bool, router ziface.IRouter) ziface.IConnection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool),
		msgChan:make(chan []byte), //msgChan初始化
	}

}