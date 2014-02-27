package service

import (
	//"sync"
	//	"monitor/cache"
	"monitor/common"
	"monitor/config"
	"monitor/service/http"
	"monitor/service/online"
	"monitor/service/rpc"
	"monitor/service/tcp"
	"time"
	"uway/log"
)

var (
	changes             = make(chan *common.TimeMessage, 10)
	PushDelay float64   = 60
	delayTime time.Time = time.Now().Truncate(time.Second)
)

//启用服务
func Start(conf *config.Config) {
	if conf.Tcp.Enable {
		tcp.StartTcpService(conf.Tcp.Host, conf.Tcp.Port)
	}
	if conf.Http.Enable {
		http.StartHttpService(conf.Http.Host, conf.Http.Port)
	}
	PushDelay = conf.PushDelay
	rpc.ServerIP = conf.RPC.Host
	rpc.Port = conf.RPC.Port
	go Run()
	log.Info("service is starting...")
}

func Run() {
	go timeTick()
	for tm := range changes {
		if len(tm.Key) != 0 && online.GetCurrentTime().Before(tm.Time) { //push time
			online.SetCurrentTime(tm.Time)
		}

		//两种情况，实时推送（pushdelay==0 ,tm是由数据推送的，key不为空）；定时推送（pushdelay>0 并超时，不考虑key）；
		if PushDelay == 0 && len(tm.Key) != 0 || PushDelay > 0 && time.Since(delayTime).Seconds() >= PushDelay {
			delayTime = time.Now()
			tm.Time = online.GetCurrentTime()
			//log.Infof("now push data @ ", tm.Key, tm.Time)
			online.SendTimeMessageToAll(tm)
		}
	}
}
func Close() {
	close(changes)
}

//定时产生时间点
func timeTick() {
	times := time.Tick(time.Second * 10)
	for common.IsRun {
		<-times
		tm := &common.TimeMessage{}
		changes <- tm
	}
}

//通知时间点变化
func MessageTimeChange(t *common.TimeMessage) {
	//log.Info("time change", t.Time)
	changes <- t
}
