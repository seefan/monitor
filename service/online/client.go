package online

import (
	"encoding/json"
	"monitor/common"
	"monitor/log"
	"monitor/service/cell"
	"sync"
	"time"
)

//client
type Client struct {
	Name           string
	OutputChan     chan []byte
	IsRun          bool      //是否正在运行
	IsLogin        bool      //是否已登陆
	IsPlayBack     bool      //是否是处于回放状态
	PlayBackTime   time.Time //回放时下发数据时间点
	dataParams     map[string]*common.DataParam
	lock           sync.RWMutex
	UUID           string
	LastUpdateTime time.Time
	Key            string
}

func (this *Client) UpdateTime() {
	this.LastUpdateTime = time.Now()
	FreshSession(this.UUID)
}

//客户端关闭接口
type ClientClose interface {
	Close()
}

//send cmd to client
func (this *Client) Send(cmd *common.ResponseData) {
	if jscmd, err := json.Marshal(cmd); err == nil {
		this.OutputChan <- jscmd
	}
}

//read cmd from byte
func (this *Client) Read(data []byte, cmd *common.RequestData) {
	if err := json.Unmarshal(data, cmd); err != nil {
		log.Error("cmd format error,cmd string is %s,error is %s", string(data), err.Error())
	}
}

//注册一个新的推送类型
func (this *Client) AddRequest(req *common.DataParam) {
	this.lock.Lock()
	if this.dataParams == nil {
		this.dataParams = make(map[string]*common.DataParam)
	}
	this.dataParams[req.Id] = req
	this.lock.Unlock()
	//println("register request", req.Id)
}

//删除一个推送类型
func (this *Client) RemoveRequest(reqId string) {
	//println("delete request", reqId)
	this.lock.Lock()
	delete(this.dataParams, reqId)
	this.lock.Unlock()
}

//send cmd to client
func (this *Client) SendBytes(jscmd []byte) {
	this.OutputChan <- jscmd
}

//新收到一个时间点变化，收集数据并下发
func (this *Client) Processing(tm *common.TimeMessage) {
	//首先看本类数据是需要下发，查字典
	//println("get time message is ", tm.Time.String(), tm.Key)
	this.lock.RLock()
	for _, v := range this.dataParams {
		//log.Info("dataparams is ", v.DataKey == tm.Key, len(tm.Key) == 0)
		if v.DataKey == tm.Key || len(tm.Key) == 0 { //同类数据
			//log.Info("celling   is ", v.Id)
			go this.celling(v, tm.Time)
		}
	}
	this.lock.RUnlock()
}

//处理各类推送数据
func (this *Client) celling(dp *common.DataParam, t time.Time) {
	re := cell.GetData(dp, t)
	//log.Info("get dp %s is ", dp.Id, re != nil)
	if re != nil {
		out := common.ResponseData{2, 0, re, dp}
		this.Send(&out)
	}
}

//客户端关闭方法
func (this *Client) Close() {
	close(this.OutputChan)
}
