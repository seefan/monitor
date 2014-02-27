package cell

import (
	"monitor/cache"
	"monitor/common"
	//	"monitor/log"
	"monitor/common/compute"
	"sort"
	"time"
)

var (
//Config *config.Config
)

//单时间点，单指标，多网元
//主要查多个网元指定时间的排序
func GetOrderData(dp *common.DataParam, t time.Time) *common.OutputParam {
	tmp := []*common.DataRow{}
	period := 1
	if p, ok := cache.Get(cache.FormatKey("System", "Period", dp.DataKey)); ok {
		period = compute.AsInt(p)
	}
	//	limit := []string{}
	if len(dp.LimitNodeId) > 0 { //双重限制，与LimitId同时起作用

	} else {
		for _, v := range dp.LimitId { //取出所有的数据
			//			key := cache.FormatKey(t.Format(common.GetTimeFormat(period)), dp.DataKey, v)
			if r, ok := cache.GetDataRow(dp.DataKey, t.Format(common.GetTimeFormat(period)), v); ok { //都是最新一个时间点的，都在内存
				tmp = append(tmp, r)
			}
		}
	}
	//排序
	sort.Sort(common.ByKey{tmp, dp.OrderKey})
	//限制输出
	result := new(common.OutputParam)
	if len(tmp) > 0 {
		result.RowMap = make(map[string]int)
		for i, k := range dp.OutputKey {
			result.RowMap[k] = i
		}
		result.Time = time.Now()
	}
	if dp.OrderBy == common.Desc {
		for i := len(tmp); i > 0 && i > len(tmp)-dp.DataRange; i-- {
			result.Rows = append(result.Rows, createRow(dp, tmp[i-1]))
		}
	} else {
		for i := 0; i < len(tmp) && i < dp.DataRange; i++ {
			result.Rows = append(result.Rows, createRow(dp, tmp[i]))
		}
	}
	return result
}
