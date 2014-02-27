package cache

import (
	"fmt"
	"monitor/common"
	"monitor/config"
	"monitor/log"
	//	"sync"
)

var (
	//	lock        sync.RWMutex
	TimeChanged func(*common.TimeMessage)
	staticCache = make(map[string]interface{})
	timeCache   = make(map[string]*TimeCache)
)

//取缓存的数据行，如果在内部缓存中没有就到外部缓存找。前提是时间已过期
func GetDataRow(dataKey string, timeKey string, id string) (*common.DataRow, bool) {
	//lock.RLock()
	//defer lock.RUnlock()
	if data, ok := timeCache[dataKey]; ok {
		if v, ok := data.Get(timeKey, id); ok {
			//println("neibu cache")
			return v.(*common.DataRow), true
		} else { //在内部缓存没有找到，在外部缓存查找
			//println("waibu cache")
			if !data.ExistsTime(timeKey) { //内部时间已过期
				//println("waibu store cache")
				return getFromStore(FormatStoreKey(dataKey, id), timeKey)
			}
		}
	}
	return nil, false
}

//get key value
func Get(key string) (interface{}, bool) {
	if v, ok := staticCache[key]; ok {
		return v, true
	}
	return nil, false
}
func Set(key string, value interface{}) {
	staticCache[key] = value
}

//按固定格式连接关键字，组成key
func FormatKey(inKeyType string, inkeyValue interface{}, outKeyName string) string {
	return fmt.Sprintf("%s:%v:%s", inKeyType, inkeyValue, outKeyName)
}

//外面用key
func FormatStoreKey(inKeyType string, inkeyValue interface{}) string {
	return fmt.Sprintf("%s:%v", inKeyType, inkeyValue)
}

func CloseCache() {
	for _, v := range timeCache {
		v.Close()
	}
}

//增加一条记录
func AddRowsToCache(nodeName string, rows *common.DataRows) {
	//log.Infof("receive %s--%s", nodeName, rows.Time)
	if c, ok := timeCache[nodeName]; ok {
		c.Add(rows)
	}
	//log.Infof("end receive %s--%s", nodeName, rows.Time)
}

//初始化数据关系
func initRelation(cfg *config.RelationConfig) {
	for _, r := range cfg.Relations {
		Set(FormatKey("System", "Relation", r.Id), r.Childs)
	}
}

//初始化缓存
func Init(cfg *config.Config, rcfg *config.RelationConfig) {
	log.Info("cache initing")
	initRedis(cfg.Redises)
	initTimeCache(cfg)
	if rcfg != nil {
		initRelation(rcfg)
	}
}
func addTimeCache(input *config.Input) {
	if input.Enable && input.EnableStore {
		c := NewTimeCache(input.Id, input.PrimaryKey.Value)
		timeCache[input.Id] = c
	}
}
