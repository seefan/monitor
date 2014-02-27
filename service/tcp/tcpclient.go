package tcp

import (
	"monitor/log"
	"monitor/service/online"
	"net"
)

//tcp client
type TcpClient struct {
	online.Client
	Connect *net.TCPConn
}

func (this *TcpClient) Close() {
	this.Client.Close()
	if this.Connect != nil {
		if err := this.Connect.Close(); err != nil {
			log.Errorf("client %s close with error:%s", this.Name, err.Error())
		}
	}
}
