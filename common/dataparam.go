package common

import (
	"encoding/json"
	"monitor/common/compute"
	"monitor/log"
	"time"
)

const (
	Desc = iota
	Asc
)

func (this *RequestData) GetString(key string) string {
	if v, ok := this.Param[key]; ok {
		return compute.AsString(v)
	}
	return ""
}
func (this *RequestData) GetStringArray(key string) []string {
	if v, ok := this.Param[key]; ok {
		if v, ok := v.([]interface{}); ok {
			re := []string{}
			for _, s := range v {
				re = append(re, compute.AsString(s))
			}
			return re
		} else {
			log.Infof("param %s value error", key, v)
		}
	}
	return nil
}
func (this *RequestData) GetFloat(key string) float64 {
	if v, ok := this.Param[key]; ok {
		compute.AsFloat64(v)
	}
	return -1
}
func (this *RequestData) GetInt64(key string) int64 {
	if v, ok := this.Param[key]; ok {
		return compute.AsInt64(v)
	}
	return -1
}
func (this *RequestData) GetInt(key string) int {
	if v, ok := this.Param[key]; ok {
		return compute.AsInt(v)
	}
	return -1
}

//input data
type RequestData struct {
	Pid   int
	Param map[string]interface{}
}

// response data
type ResponseData struct {
	Pid    int         //协议号
	Status int         //状态
	Output interface{} //输出内容
	Param  interface{}
}

//数据推送请求参数
type DataParam struct {
	FillKey     string   //回填from的key，不指定就不进行回填
	TimeRange   int      //时间点个数
	DataRange   int      //最多显示多少条数据，多时间点时，此项无效
	DataKey     string   //取哪类数据，与数据节点id对应
	LimitId     []string //限制只取哪几个id的数据
	LimitNodeId string   //限制父节点id，与LimitId互斥	 [排行榜专用]
	OrderKey    string   //排序指标		 [排行榜专用]
	OrderBy     int      //排序方式，		 [排行榜专用]
	OutputKey   []string //输出的字段
	Id          string   //请求ID，在客户端要区分开，不同的请求使用不同的id，在不需要时要删除该推送请求
}

func (this *DataParam) ToString() string {
	if s, err := json.Marshal(this); err == nil {
		return string(s)
	} else {
		return err.Error()
	}

}

//数据输出结构
type OutputParam struct {
	Time   time.Time
	RowMap map[string]int
	Rows   [][]interface{}
}

//func (this *DataParam) Hashcode() int32 {
//	tostring := fmt.Sprint(this.GetType, this.DataRange, this.DataType, strings.Join(this.LimitId, ":"), this.PrimaryKey, this.Order, strings.Join(this.OutputKey, ":"), this.Id)
//	return hash.HashCode(tostring)
//}

//用于提示时间消息的结构
type TimeMessage struct {
	Key  string
	Time time.Time
}
