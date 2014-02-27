package http

import (
	"code.google.com/p/go.net/websocket"
	"monitor/log"
	"monitor/service/online"
)

//http client
type HttpClient struct {
	online.Client
	Connect *websocket.Conn
}

func (this *HttpClient) Close() {
	this.Client.Close()
	if this.Connect != nil {
		if err := this.Connect.Close(); err != nil {
			log.Errorf("client %s close with error:%s", this.Name, err.Error())
		}
	}
}
