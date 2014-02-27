package cell

import (
	"monitor/cache"
	"monitor/common"
	"monitor/common/compute"
	//	"monitor/log"
	"time"
)

//单时间点，多指标，单网元
//取一个网元指定时间的多指标，主要用于查询某网元的详细指标
func GetColumnsData(dp *common.DataParam, t time.Time) *common.OutputParam {
	tmp := []*common.DataRow{}
	period := 1
	if p, ok := cache.Get(cache.FormatKey("System", "Period", dp.DataKey)); ok {
		period = compute.AsInt(p)
	}
	//key := cache.FormatKey(t.Format(common.GetTimeFormat(period)), dp.DataKey, dp.LimitId[0])
	//log.Info("id=%s", key)
	if r, ok := cache.GetDataRow(dp.DataKey, t.Format(common.GetTimeFormat(period)), dp.LimitId[0]); ok {
		tmp = append(tmp, r)
	}

	//限制输出
	result := new(common.OutputParam)
	if len(tmp) > 0 {
		result.RowMap = make(map[string]int)
		for i, k := range dp.OutputKey {
			result.RowMap[k] = i
		}
		result.Time = time.Now()
		result.Rows = append(result.Rows, createRow(dp, tmp[0]))
	}
	return result
}
