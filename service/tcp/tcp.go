//每个请求会登陆一个Client，如果登陆成功就保存起来，继续运行。每个Client有两个协程，分别读取和写入。
//写入使用一个chanel，保证同步
//每个读写请求一但出错，立即抛弃Client，等待客户端再次登陆。连接会自动关闭
//
//
//

package tcp

import (
	"bufio"
	//	"encoding/json"
	"fmt"
	"net"
	//	"time"
	"code.google.com/p/go-uuid/uuid"
	"monitor/common"
	"monitor/log"
	"monitor/service/online"
	"monitor/service/protocol"
)

var (
	maxRead = 1024

//	Message = make(chan time.Time)
)

//tcp service
type TcpService struct {
	Host  string
	Port  int
	IsRun bool
}

//init tcp service
func (this *TcpService) init(host string, port int) {
	this.Host = host
	this.Port = port
}

//create new http service
//with host and port
func StartTcpService(host string, port int) {
	h := new(TcpService)
	h.init(host, port)
	go h.Run()

}

//tcp service start
func (this *TcpService) Run() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", this.Host, this.Port))
	if err != nil {
		panic("resolve  tcp address error")
	}
	link, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic("start tcp service error")
	}
	defer link.Close()
	this.IsRun = true
	log.Infof("tcp service started at %s:%d", this.Host, this.Port)
	for {
		conn, err := link.AcceptTCP()
		if err != nil {
			log.Error(err)
			continue
		}

		connFrom := conn.RemoteAddr().String()
		log.Infof("accept new client from %s\n", connFrom)
		c := new(TcpClient)
		c.Connect = conn
		c.UUID = uuid.New()
		c.IsRun = true
		c.OutputChan = make(chan []byte, 10)
		go this.HandleRequest(c)
		go this.HandleResponse(c)
	}
}

//send cmd to client
func (this *TcpService) HandleResponse(c *TcpClient) {
	defer online.Delete(c.Name)
	c.Connect.SetWriteBuffer(common.BufferSize)
	//c.Connect.SetNoDelay(true)
	var err error
	for cmd := range c.OutputChan {
		//cmd := <-c.OutputChan
		cmd = append(cmd, '\n')
		if _, err = c.Connect.Write(cmd); err != nil {
			log.Error("send cmd error and login out,%s\n", err.Error())
			c.IsRun = false
			break //结束协程
		}
	}
}

//handle tcp request
func (this *TcpService) HandleRequest(c *TcpClient) {
	defer online.Delete(c.Name)
	nr := bufio.NewReader(c.Connect)
	var err error
	var data []byte
	cmd := new(common.RequestData)
	for this.IsRun && c.IsRun {
		data, err = nr.ReadBytes(common.NewLine) //读取一行数据
		if err != nil {
			c.IsRun = false
			log.Error("tcp service read error", err.Error())
			break
		}
		c.Read(data, cmd)
		//if err := json.Unmarshal(data, &cmd); err != nil {
		//	log.Error("cmd format error,", err.Error())
		//}
		//log.Info("server cmd", cmd)
		protocol.Exec(cmd, &c.Client)
	}

}
