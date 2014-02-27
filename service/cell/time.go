package cell

import (
	"monitor/cache"
	"monitor/common"
	"monitor/common/compute"
	//	"monitor/log"
	"sort"
	"time"
)

//时间线，单指标，单网元
//查询一段时间内指定网元，某指标的走势
func GetTimeLineData(dp *common.DataParam, t time.Time) *common.OutputParam {
	tmp := []*common.DataRow{}
	//log.Info("timeline @", dp.Id, t)
	period := 1
	if p, ok := cache.Get(cache.FormatKey("System", "Period", dp.DataKey)); ok {
		period = compute.AsInt(p)
	}
	for i := 0; i < dp.TimeRange; i++ { //取出所有的数据
		tmpTime := t.Add(time.Minute * time.Duration(-i))
		//		key := cache.FormatKey(tmpTime.Format(common.GetTimeFormat(period)), dp.DataKey, dp.LimitId[0])
		//log.Info("search key is ", key)
		if r, ok := cache.GetDataRow(dp.DataKey, tmpTime.Format(common.GetTimeFormat(period)), dp.LimitId[0]); ok {
			tmp = append(tmp, r)
			//log.Infoln("find cache", r.Row)
		}
	}
	//log.Info("result1")
	//排序
	sort.Sort(common.ByTime{tmp})
	//log.Info("result2")
	result := new(common.OutputParam)
	if len(tmp) > 0 {
		result.RowMap = make(map[string]int)
		for i, k := range dp.OutputKey {
			result.RowMap[k] = i
		}
		result.Time = t
	}
	//log.Info("result3")
	//限制输出
	//start := 0
	//if len(tmp) >= dp.DataRange {
	//	start = len(tmp) - dp.DataRange
	//}
	//for i := start; i < len(tmp); i++ {
	//	result.Rows = append(result.Rows, createRow(dp, tmp[i]))
	//}
	for _, r := range tmp {
		result.Rows = append(result.Rows, createRow(dp, r))
	}
	//log.Info("result")
	return result
}
