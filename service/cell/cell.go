package cell

import (
	"monitor/cache"
	"monitor/common"
	"monitor/log"
	"time"
)

//处理单元入口
func GetData(dp *common.DataParam, t time.Time) *common.OutputParam {
	if dp.TimeRange == 1 && len(dp.LimitId) == 1 { //单网元，单时间点，用于指标详细
		log.Infoln("开始处理过程，单网元，单时间点，用于指标详细")
		return GetColumnsData(dp, t)
	} else if dp.TimeRange > 1 && len(dp.LimitId) == 1 { //单网元，多时间点，用于走势图
		log.Infoln("开始处理过程，单网元，多时间点，用于走势图")
		return GetTimeLineData(dp, t)
	} else if dp.TimeRange == 1 && (len(dp.LimitId) > 1 || len(dp.LimitNodeId) > 0) { //单时间点，多网元，用于排行榜
		log.Infoln("开始处理过程，单时间点，多网元，用于排行榜")
		return GetOrderData(dp, t)
	} else { //多时间点，多网元，暂不支持
		log.Error("开始处理过程，不被支持的过程", dp.ToString())
		return nil
	}
}

//从内部存贮创建一个新行
func createRow(dp *common.DataParam, r *common.DataRow) []interface{} {
	var row []interface{}
	for _, k := range dp.OutputKey {
		if v, ok := r.GetValue(k); ok {
			row = append(row, v)
		} else { //回填数据，从缓存中按回填key取值,key值在dp中指定
			//log.Infoln("key is ", dp.FillKey)
			v = getFillValue(dp.FillKey, r, k)
			row = append(row, v)
		}
	}
	return row
}

//回填
func getFillValue(fillKey string, r *common.DataRow, k string) interface{} {
	if len(fillKey) > 0 {
		key := cache.FormatKey(fillKey, r.GetKey(fillKey), k)
		//	log.Infoln("key is ", key)
		if tv, ok := cache.Get(key); ok {
			return tv
		}
	}
	return nil
}
func getStoreCache() *common.DataRow {
	return nil
}
