/*
用cookie保持连接，每10分钟换一个auth值。
与tcp服务不同，websocket的连接断开后不需要重新登陆，只需要验证Auth的登陆信息即可。
不同的单元可以保持不同的连接，简化客户端的代码
不同的连接产生不同的Client，同一个连接内不允许多次登陆。多次登陆也不会重新生成Client
登陆后生成key，用于同一用户多连接,and reset key
*/
package http

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"monitor/common"
	"monitor/log"
	"monitor/service/online"
	"monitor/service/protocol"
	"net/http"
	"time"
)

var (
	session = make(map[string]string) //在线用户
)

const (
	Sec_Auth_Uid = "Sec_Auth_Uid"
	Sec_Auth_Key = "Sec_Auth_Key"
)

//http service
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

//处理单独的Http请求
func (this *HttpService) HandleHttp(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	uid := req.Header.Get(Sec_Auth_Uid)
	key := req.Header.Get(Sec_Auth_Key)
	c := this.getClient(uid, key)
	defer c.Close()

	if data := req.FormValue("json"); len(data) > 0 {
		cmd := new(common.RequestData)
		c.Read([]byte(data), cmd)
		protocol.Exec(cmd, &c.Client)
	} else {
		cmd := common.ResponseData{0, -1, "json string is empty", ""}
		c.Send(&cmd)
	}
	for {
		timeout := time.After(time.Nanosecond * 10) //10ns就超时
		select {
		case data := <-c.OutputChan:
			if _, err := w.Write(data); err != nil {
				log.Error("%s send cmd error and login out,%s\n", c.Name, err.Error())
			}
			break
		case <-timeout:
			goto end
		}
	}
end:
}

//按cookei取在线用户
func (this *HttpService) getClient(uid string, key string) *HttpClient {
	c := new(HttpClient)
	if s, _, ok := online.GetSession(uid); ok && s == key { //key相同时认为已登陆过的
		if tc, ok := online.Get(uid); ok { //第一个登陆的Client的uuid被当作了uid
			c.Name = tc.Name
			c.IsLogin = tc.IsLogin
			c.IsPlayBack = c.IsPlayBack
		}
	}
	if c == nil {
		c = new(HttpClient)
	}
	c.UUID = uuid.New()
	c.IsRun = true
	c.OutputChan = make(chan []byte, 10)
	return c
}

////设置cookie
//func (this *HttpService) setCookie(c *HttpClient, w http.ResponseWriter, age int) {
//	id := &http.Cookie{Name: "monitor_http_key", Value: common.HashString(c.Name + ":monitor"), Path: "/", MaxAge: age, HttpOnly: true, Secure: true}
//	name := &http.Cookie{Name: "monitor_http_name", Value: c.Name, Path: "/", MaxAge: age, HttpOnly: true, Secure: true}
//	http.SetCookie(w, id)
//	http.SetCookie(w, name)
//}

//处理Socket请求
func (this *HttpService) HandleSocket(ws *websocket.Conn) {
	connFrom := ws.RemoteAddr()
	log.Infof("accept new http client from %s\n", connFrom)
	uid := ws.Request().Header.Get(Sec_Auth_Uid)
	key := ws.Request().Header.Get(Sec_Auth_Key)
	c := this.getClient(uid, key)
	c.Connect = ws
	go this.HandleResponse(c)
	this.HandleRequest(c)
}

//主运行方法
//根目录为标准Http请求
//socket为Socket请求
func (this *HttpService) Run() {
	this.IsRun = true
	http.Handle("/", http.FileServer(http.Dir("."))) // <-- note this line
	http.HandleFunc("/http", this.HandleHttp)
	http.Handle("/socket", websocket.Handler(this.HandleSocket))
	log.Infof("http service started at %s:%d", this.Host, this.Port)
	log.Info("http socket host on /http and /socket ")
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
	var data []byte
	for this.IsRun && c.IsRun {
		if err = websocket.Message.Receive(c.Connect, &data); err != nil {
			log.Infof("%s can't received cmd", c.Name, err.Error())
			c.IsRun = false
			break
		}
		c.Read(data, cmd)
		protocol.Exec(cmd, &c.Client)
		c.UpdateTime()
	}
}
