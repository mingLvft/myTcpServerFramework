package main

import (
	"fmt"
	"io"
	"myServerDemo/zinx/znet"
	"net"
	"time"
)

/**
 * 模拟客户端
 */
func main() {
	fmt.Println("client0 start...")

	time.Sleep(1 * time.Second)

	// 1 直接链接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的message消息 MsgID: 0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		// 服务器就应该给我们回复一个message数据，MsgID：1 ping ping ping

		// 先读取流中的head部分 得到ID和Datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error ", err)
			break
		}

		// 将二进制的head拆包到msg 结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 根据Datalen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error , ", err)
				return
			}

			fmt.Println("————> Recv Server Msg : ID = ", msg.Id, "len = ", msg.DataLen, ", data = ", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
