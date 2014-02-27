/*
用cookie保持连接，每10分钟换一个cookie的值。
与tcp服务不同，websocket的连接断开后不需要重新登陆，只需要验证cookie的登陆信息即可。
不同的单元可以保持不同的连接，简化客户端的代码
*/
package http

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	//	"encoding/json"
	"fmt"
	"monitor/common"
	"monitor/log"
	"monitor/service/online"
	"monitor/service/protocol"
	"net/http"
)

type HttpService struct {
	Host  string
	Port  int
	IsRun bool
}

//init tcp service
func (this *HttpService) init(host string, port int) {
	this.Host = host
	this.Port = port
	this.IsRun = true
}

//create new http service
//with host and port
func StartHttpService(host string, port int) {
	h := new(HttpService)
	h.init(host, port)
	go h.Run()

}

//处理Http请求
func (this *HttpService) HandleHttp(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	c := this.getClient(req)
	defer c.Close()
	log.Info("login0", c.Name, c.IsLogin)
	if data := req.FormValue("json"); len(data) > 0 {
		cmd := new(common.RequestData)
		c.Read([]byte(data), cmd)
		log.Info("login1", c.Name, c.IsLogin)
		protocol.Exec(cmd, &c.Client)
		log.Info("login2", c.Name, c.IsLogin)
		if cmd.Pid == 1 && c.IsLogin {
			log.Info("client %s is login,set cookie", c.Name)
			this.setCookie(c, w, 0)
		}
		if cmd.Pid == -1 {
			this.setCookie(c, w, -1)
		}
	} else {
		cmd := common.ResponseData{0, -1, " json string is empty", ""}
		c.Send(&cmd)
	}
	data := <-c.OutputChan
	if _, err := w.Write(data); err != nil {
		log.Error("%s send cmd error and login out,%s\n", c.Name, err.Error())
	}

}

//按cookei取在线用户
func (this *HttpService) getClient(req *http.Request) *HttpClient {
	c := new(HttpClient)
	if id, err := req.Cookie("monitor_http_key"); err == nil {
		log.Info("monitor_http_key=", id.Value)
		if name, err := req.Cookie("monitor_http_name"); err == nil {
			log.Info("monitor_http_name", name.Value)
			c.Name = name.Value
			c.IsLogin = id.Value == common.HashString(c.Name+":monitor")
		}
	}
	c.UUID = uuid.New()
	c.IsRun = true
	c.OutputChan = make(chan []byte, 10)
	return c
}

//设置cookie
func (this *HttpService) setCookie(c *HttpClient, w http.ResponseWriter, age int) {
	id := &http.Cookie{Name: "monitor_http_key", Value: common.HashString(c.Name + ":monitor"), Path: "/", MaxAge: age, HttpOnly: true, Secure: true}
	name := &http.Cookie{Name: "monitor_http_name", Value: c.Name, Path: "/", MaxAge: age, HttpOnly: true, Secure: true}
	http.SetCookie(w, id)
	http.SetCookie(w, name)
}

//处理Socket请求
func (this *HttpService) HandleSocket(ws *websocket.Conn) {
	connFrom := ws.RemoteAddr()
	log.Infof("accept new http client from %s\n", connFrom)
	c := this.getClient(ws.Request())
	c.Connect = ws

	go this.HandleResponse(c)
	this.HandleRequest(c)
}

//主运行方法
//根目录为标准Http请求
//socket为Socket请求
func (this *HttpService) Run() {
	http.Handle("/", http.FileServer(http.Dir("."))) // <-- note this line
	http.HandleFunc("/http", this.HandleHttp)
	http.Handle("/socket", websocket.Handler(this.HandleSocket))
	this.IsRun = true
	log.Infof("http service started at %s:%d", this.Host, this.Port)
	log.Info("http socket host on /socket")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", this.Host, this.Port), nil); err != nil {
		log.Error("http service start error", err.Error())
	}
}

//send cmd to client
func (this *HttpService) HandleResponse(c *HttpClient) {
	var err error
	defer online.Delete(c.UUID)
	for cmd := range c.OutputChan {
		//cmd := <-c.OutputChan
		cmd = append(cmd, '\n')
		if err = websocket.Message.Send(c.Connect, string(cmd)); err != nil {
			log.Error("%s send cmd error and login out,%s\n", c.Name, err.Error())
			c.IsRun = false
			break
		}
		c.UpdateTime()
	}
}

//handle tcp request
func (this *HttpService) HandleRequest(c *HttpClient) {
	var err error
	defer online.Delete(c.UUID)
	cmd := new(common.RequestData)
	online.Set(&c.Client)
	var data []byte
	for this.IsRun && c.IsRun {
		if err = websocket.Message.Receive(c.Connect, &data); err != nil {
			log.Infof("%s can't received cmd", c.Name, err.Error())
			c.IsRun = false
			break
		}
		c.Read(data, cmd)
		protocol.Exec(cmd, &c.Client)
		if cmd.Pid == 1 && c.IsLogin {
			log.Info("client %s is login,set cookie", c.Name)
			//this.setCookie(c, c.Connect. 0)
		}
		c.UpdateTime()
	}
}
