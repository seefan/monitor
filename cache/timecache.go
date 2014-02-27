/*
针对不同的数据类型，采用不同的缓存方式
1、普通不过期缓存，使用Set和Get方法统计处理，三段式Key值
2、定时过期数据，按StoreMinute进行存贮，过期清理
*/
package cache

import (
	"monitor/common"
	"monitor/common/compute"
	"monitor/config"
	"monitor/log"
	"sync"
	"time"
)

type TimeCache struct {
	Id      string
	KeyName string                    //数据列字段名
	input   chan *common.DataRows     //数据输入窗口
	caches  map[string]*TimeCacheItem //所有的缓存数据
	lock    sync.RWMutex
}
type TimeCacheItem struct {
	Cache       map[string]interface{}
	ExpiredTime time.Time
}

func (this *TimeCache) init() {
	this.input = make(chan *common.DataRows, 100)
	this.caches = make(map[string]*TimeCacheItem)
}
func (this *TimeCache) ExistsTime(timeKey string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.caches[timeKey]
	return ok
}
func (this *TimeCache) Get(timeKey string, id string) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if item, ok := this.caches[timeKey]; ok {
		if value, ok := item.Cache[id]; ok {
			return value, true
		}
	}
	return nil, false
}

//创建一个新的缓存实例
func NewTimeCache(id string, keyName string) *TimeCache {
	c := &TimeCache{Id: id, KeyName: keyName}
	c.Start()
	return c
}

//启动实例
func (this *TimeCache) Start() {
	this.init()
	go this.Run()
	go this.Clean(30)
}

//清理过期数据
func (this *TimeCache) Clean(s time.Duration) {
	tc := time.Tick(time.Second * s)
	for t := range tc {
		ts := []string{}
		for k, v := range this.caches {
			if v.ExpiredTime.Before(t) {
				//	log.Infof("data %s is expired %v,now is %s", k, v.ExpiredTime, t.String())
				ts = append(ts, k)
			}
		}
		this.lock.Lock()
		for _, k := range ts {
			delete(this.caches, k)
		}
		this.lock.Unlock()
	}
}

//close it
func (this *TimeCache) Close() {
	close(this.input)
}

//增加一条数据
func (this *TimeCache) Add(rows *common.DataRows) {
	//	log.Info("receive new rows@", rows.Time)
	this.input <- rows
}

//运行计算
func (this *TimeCache) Run() {
	log.Infoln("cache is started", this.Id)
	period := 1
	if p, ok := Get(FormatKey("System", "Period", this.Id)); ok {
		period = compute.AsInt(p)
	}
	timeFormat := common.GetTimeFormat(period)

	for rs := range this.input { //接收数据
		//log.Infof("time changed ", rs.Time)
		currTime, err := time.Parse(timeFormat, rs.Time)
		if err != nil {
			log.Error("cache time format is error", err.Error())
			currTime = time.Now()
		}
		cacheItem, ok := this.caches[rs.Time] //这个时间点是否已记录
		if !ok {
			cacheItem = &TimeCacheItem{Cache: make(map[string]interface{})}
			this.caches[rs.Time] = cacheItem
			cacheItem.ExpiredTime = time.Now().Add(time.Minute * common.StoreMinute)
		}

		for _, vs := range rs.Rows { //分析每一条数据
			row := rs.CreateDataRow(vs)
			key := row.GetKey(this.KeyName)
			//log.Infoln("add to cache", key)
			this.lock.Lock()
			if r, ok := cacheItem.Cache[key]; ok { //如果同类数据已存在，就进行数据合并
				row.Merge(r.(*common.DataRow).Row)
			}
			cacheItem.Cache[key] = row
			this.lock.Unlock()
			//log.Infoln(key, row.Row)
			//	log.Info("key=%s", key)
			go saveToStore(FormatStoreKey(this.Id, key), row, period*common.StoreCount) //保存到外部存贮，并保存粒度*保存个数分钟
		}

		//通知外到有新数据到来
		if TimeChanged != nil {
			TimeChanged(&common.TimeMessage{Key: this.Id, Time: currTime})
		}
	}
}

func CountCache() int {
	re := 0
	for _, v := range timeCache {
		for _, cv := range v.caches {
			re += len(cv.Cache)
		}
	}
	return re
}

//初始化时间缓存节点
func initTimeCache(cfg *config.Config) {
	for _, input := range cfg.Inputs.CsvInputs {
		addTimeCache(&input.Input)
	}
	for _, input := range cfg.Inputs.RedisInputs {
		addTimeCache(&input.Input)
	}
	for _, input := range cfg.Inputs.TcpInputs {
		addTimeCache(&input.Input)
	}
	for _, input := range cfg.Inputs.TestInputs {
		addTimeCache(&input.Input)
	}
	for _, input := range cfg.Summarys {
		c := NewTimeCache(input.Id, input.Relation.PrimaryKey)
		timeCache[input.Id] = c
	}
}
