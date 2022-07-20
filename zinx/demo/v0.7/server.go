package main

import (
	"fmt"
	"io"
	"net"
	"zinx/zinx/znet"
)

// -------------------------------------------
// @file          : server.go
// @author        : binshow
// @time          : 2022/7/20 9:47 AM
// @description   :
// -------------------------------------------


func main() {

	//1. 创建一个 tcp server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listen err : ", err)
		return
	}

	//2. 开启goroutine， 负责从client的goroutine 读取数据，并解决粘包问题
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err :" , err)
		}

		go func(conn net.Conn) {
			// 创建封包和粘包对象
			dp := znet.NewDataPack()

			for true {
				// 从 conn 中读取 header 中的数据
				headData := make([]byte , dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head err :" , err)
					break
				}

				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err:" , err)
					return
				}

				if msgHead.GetDataLen() > 0 {

					// msg 中有 data数据，再次读取conn
					msg := msgHead.(*znet.Message)
					msg.Data = make([]byte , msg.GetDataLen())

					_ , err := io.ReadFull(conn , msg.Data)
					if err != nil {
						fmt.Println("server unpack data err :" , err)
						return
					}
					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)


	}


}