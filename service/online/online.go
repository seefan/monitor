package online

import (
	"encoding/json"
	"monitor/common"
	"monitor/log"
	"sync"
	"time"
)

var (
	onlineClient = make(map[string]*Client)
	currentTime  time.Time //time from data push
	lock         sync.RWMutex
)

//通知所有客户端有时间消息
func SendTimeMessageToAll(tm *common.TimeMessage) {
	log.Infoln("send message is ", tm.Time, tm.Key, len(onlineClient))
	for _, c := range onlineClient {
		go c.Processing(tm)
	}
}

//取得当前时间
func GetCurrentTime() time.Time {
	lock.RLock()
	defer lock.RUnlock()
	return currentTime
}

//设置当前时间
func SetCurrentTime(t time.Time) {
	lock.Lock()
	defer lock.Unlock()
	currentTime = t
}

//send cmd to all client
func SendToAll(cmd interface{}) {
	if jscmd, err := json.Marshal(cmd); err == nil {
		for _, c := range onlineClient {
			c.SendBytes(jscmd)
		}
	}
}

//将客户端加到在线列表
func Set(c *Client) {
	if _, ok := onlineClient[c.UUID]; !ok { //说明已存在同名用户
		onlineClient[c.UUID] = c
	}

}

//取指定客户端
func Get(name string) (client *Client, ok bool) {
	if v, exists := onlineClient[name]; exists {
		client = v
		ok = true
	}
	return nil, false
}

//将客户端从在线列表中删除
func Delete(name string) {
	if c, ok := onlineClient[name]; ok {
		var cc ClientClose
		cc = c
		cc.Close()
		delete(onlineClient, name)
		log.Infof("client: %s close and remove from online", c.Name)
	}
}
