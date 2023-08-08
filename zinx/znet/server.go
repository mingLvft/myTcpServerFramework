package znet

import (
	"errors"
	"fmt"
	"myServerDemo/zinx/ziface"
	"net"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
}

// 定义当前客户端链接的所绑定的handlerapi（目前这个handler是写死的，以后应该由用户自定义handler方法）
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显的业务
	fmt.Println("[Conn Handler] CallBackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server success, ", s.Name, " success, Listenning...")
		var cid uint32
		cid = 0

		// 3 阻塞的等待客户端连接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP")
				continue
			}

			//// 已经与客户端建立链接，做一些业务，做一个最基本的最大512字节长度的回显业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("receive buf err", err)
			//			continue
			//		}
			//
			//		fmt.Printf("receive client buf %s, cnt %d\n", buf, cnt)
			//		// 回显功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err ", err)
			//			continue
			//		}
			//	}
			//}()

			// 将处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

/**
 * 初始化Server模块的方法
 */
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
