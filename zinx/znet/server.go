package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/zinx/ziface"
)

// -------------------------------------------
// @file          : server.go
// @author        : binshow
// @time          : 2022/7/17 11:52 AM
// @description   : 服务接口的具体实现
// -------------------------------------------


//Server 服务端实现类
type Server struct {
	Name   		string
	IPVersion 	string
	IP			string	// 监听的ip地址
	Port        int		// 端口号

	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	//Router 		ziface.IRouter
	MsgHandle  		ziface.IMsgHandle
}

func NewServer(name string) ziface.IServer {
	//先初始化全局配置文件
	utils.GlobalObject.Reload()

	s:= &Server {
		Name :utils.GlobalObject.Name,//从全局参数获取
		IPVersion:"tcp4",
		IP:utils.GlobalObject.Host,//从全局参数获取
		Port:utils.GlobalObject.TcpPort,//从全局参数获取
		MsgHandle: NewMsgHandle(),
	}
	return s
}

// 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	// Listen
	go func() {
		//1. get a tcp address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}


		//2. listen the tcp address
		listener, err := net.ListenTCP("tcp" , addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}


		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")
		var cid uint32


		//3. start to accept net conn
		for true {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 TODO Server.Start() 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的

			dealConn := NewConnection(conn, cid, false, s.MsgHandle)
			cid++
			go dealConn.Start()	// 启动当前的链接任务
		}
	}()

}

// 暂时写死，后面可以让用户自定义来实现
// 定义当前客户端链接所绑定的 handle api
func callbackToClient(conn *net.TCPConn , data []byte , cnt int) error {
	fmt.Println("[conn handle] callback to client")
	if _ , err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err , " , err)
		return errors.New("callback to client error")
	}
	return nil
}


func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name " , s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}


func (s *Server) Server() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}


//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server)AddRouter(msgId uint32 , router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId , router)
	fmt.Println("Add Router succ! " )
}