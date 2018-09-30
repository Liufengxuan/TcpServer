package main

import (
	"TcpServer/handle"
	"log"
	"net"
	"runtime"
	"time"
)

func main() {
	//监听 begin
	var connChan = make(chan net.Conn, 10) //存储连接。
	var errChan = make(chan error)         //存储生产者和消费者的错误
	var isReStart = true
	var address = ip + ":" + port

	listener, err := net.Listen("tcp", address)
	defer listener.Close()
	if err != nil {
		log.Printf("[%s：地址使用失败]\n", address)
	}
	log.Printf("[准备监听 %s]\n", address)

	for {
		if isReStart {
			runtime.GOMAXPROCS(maxProcs)
			go waitConnection(listener, connChan, errChan)
			log.Println("[监听进程已启动]")
			isReStart = false
		}

		select {
		case conn := <-connChan:
			go handle.HandlerConn(conn)
		case err := <-errChan:
			if reListenNum > 0 {
				log.Printf("[主程序异常退出，将在 3 秒后重启 异常原因：%s ]\n", err)
				isReStart = true
				reListenNum--
				time.Sleep(time.Second * 3000)
				continue
			} else {
				log.Println("[主进程尝试多次启动后失败]")
				return
			}

		}
	}

}

func waitConnection(listener net.Listener, connChan chan<- net.Conn, errChan chan<- error) {

	for {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- err
			runtime.Goexit()
		}
		connChan <- conn
	}
}
