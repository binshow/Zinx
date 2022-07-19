package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
	"zinx/zinx/ziface"
)

// -------------------------------------------
// @file          : server_test.go
// @author        : binshow
// @time          : 2022/7/17 12:09 PM
// @description   : test 实现的服务端
// -------------------------------------------

// mock a client
func ClientTest() {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("hello binshow"))
		if err !=nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf(" server call back : %s, cnt = %d\n", buf,  cnt)
		time.Sleep(1*time.Second)
	}
}

func TestServer_Server(t *testing.T) {

	//1 创建一个server 句柄 s
	s := NewServer("[zinx V0.1]")

	//客户端测试
	go ClientTest()

	//2 开启服务
	s.Server()
}



//======= 测试自定义路由==========
//ping test 自定义路由
type PingRouter struct {
	BaseRouter //一定要先基础BaseRouter
}

//Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}
//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}

func TestServer_AddRouter(t *testing.T) {
	//创建一个server句柄
	s := NewServer("[zinx V0.3]")
	s.AddRouter(1 , &PingRouter{})

	//客户端测试
	go ClientTest()

	//2 开启服务
	s.Server()

}





