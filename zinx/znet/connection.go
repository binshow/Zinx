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

	//handleAPI ziface.HandleFunc	 // 当前链接所绑定的处理业务方法 API(一个连接绑定一个)

	//该连接的处理方法router ， 等价于上面的 handleAPI
	//Router  ziface.IRouter

	MsgHandle ziface.IMsgHandle

	ExitBuffChan chan bool // 告知当前链接已经退出or停止

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
		req := &Request{
			conn:c,
			//data:buf,
			msg: msg,
		}

		// 调用对应的路由去处理消息
		go c.MsgHandle.DoMsgHandler(req)

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start() , connID = " , c.ConnID)
	// 启动从当前链接的读数据的业务
	// todo 启动从当前链接写数据的业务

	go c.startReader()

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

	//写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
}



func NewConnection(conn *net.TCPConn, connID uint32, isClosed bool, msgHandle ziface.IMsgHandle) ziface.IConnection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandle:       msgHandle,
		ExitBuffChan: make(chan bool),
	}

}