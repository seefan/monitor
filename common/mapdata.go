package common

import (
	"fmt"
	"monitor/common/compute"
	"strconv"
)

//数据集，批处理用
type DataRows struct {
	Time   string
	RowMap *DataRowMap
	Rows   [][]interface{}
}
type DataRow struct {
	Time   string
	RowMap *DataRowMap
	Row    []interface{}
}

//创建一个数据行结构
func (this *DataRows) CreateDataRow(v []interface{}) *DataRow {
	return &DataRow{
		Time:   this.Time,
		RowMap: this.RowMap,
		Row:    v,
	}
}

//节点中流转的数据结构
type DataRowMap struct {
	Key         map[string]int
	Sum         map[string]int
	Max         map[string]int
	Min         map[string]int
	Count       map[string]int
	ColumnCount int //列统计
}

func (this *DataRows) CloneMap() (rowmap *DataRowMap) {
	rowmap = new(DataRowMap)
	rowmap.Key = make(map[string]int)
	rowmap.Sum = make(map[string]int)
	rowmap.Max = make(map[string]int)
	rowmap.Min = make(map[string]int)
	rowmap.Count = make(map[string]int)
	if this.RowMap == nil {
		println("column map  is not found")
		return
	}
	i := 0
	for k, v := range this.RowMap.Key {
		rowmap.Key[k] = v
		i = i + 1
	}
	for k, v := range this.RowMap.Sum {
		rowmap.Sum[k] = v
		i = i + 1
	}
	for k, v := range this.RowMap.Max {
		rowmap.Max[k] = v
		i = i + 1
	}
	for k, v := range this.RowMap.Min {
		rowmap.Min[k] = v
		i = i + 1
	}
	for k, v := range this.RowMap.Count {
		rowmap.Count[k] = v
		i = i + 1
	}
	rowmap.ColumnCount = i
	return
}

//取得指标的key值
func (this *DataRows) GetKey(key string, row []interface{}) string {
	if this.RowMap == nil {
		return fmt.Sprintf("column map  is not found")
	}
	if i, ok := this.RowMap.Key[key]; ok {
		return compute.AsString(row[i])
	}
	//log.Infoln(this.RowMap)
	return fmt.Sprintf("column %s is not found", key)
}

//取得特定的值
func (this *DataRow) GetValue(key string) (interface{}, bool) {
	if i, ok := this.RowMap.Key[key]; ok {
		return this.Row[i], true
	}
	return nil, false
}

//合并数据
func (this *DataRow) Merge(row []interface{}) {
	if this.RowMap.Sum != nil {
		this.Sum(row)
	}
	if this.RowMap.Sum != nil {
		this.Max(row)
	}
	if this.RowMap.Sum != nil {
		this.Min(row)
	}
	if this.RowMap.Sum != nil {
		this.Count(row)
	}
}

//按sum计算指标
func (this *DataRow) Sum(row []interface{}) {
	for _, i := range this.RowMap.Sum {
		this.Row[i] = compute.Add(this.Row[i], row[i])
	}
}

//按max计算指标
func (this *DataRow) Max(row []interface{}) {
	for _, i := range this.RowMap.Max {
		this.Row[i] = compute.Add(this.Row[i], row[i])
	}
}

//按min计算指标
func (this *DataRow) Min(row []interface{}) {
	for _, i := range this.RowMap.Min {
		this.Row[i] = compute.Add(this.Row[i], row[i])
	}
}

//按count计算指标
func (this *DataRow) Count(row []interface{}) {
	for _, i := range this.RowMap.Count {
		this.Row[i] = compute.Add(this.Row[i], 1)
	}
}

//取得指标的key值
func (this *DataRow) GetKey(key string) string {
	if this.RowMap == nil {
		return fmt.Sprintf("column map  is not found")
	}
	if i, ok := this.RowMap.Key[key]; ok {
		switch v := this.Row[i].(type) {
		case string:
			return v
		case []byte:
			return string(v)
		case int8, int16, int32, int64, int, uint, uint16, uint32, uint64, uint8:
			return fmt.Sprintf("%d", v)
		case float32:
			return strconv.FormatFloat(float64(v), 'f', 0, 32)
		case float64:
			return strconv.FormatFloat(v, 'f', 0, 64)
		}
		return fmt.Sprintf("%v", this.Row[i])
	}
	//log.Infoln(this.RowMap)
	return fmt.Sprintf("column %s is not found", key)
}

type DataRowSorter []*DataRow

func (this DataRowSorter) Len() int {
	return len(this)
}

func (this DataRowSorter) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type ByTime struct {
	DataRowSorter
}

func (this ByTime) Less(i, j int) bool {
	return this.DataRowSorter[i].Time < this.DataRowSorter[j].Time
}

type ByKey struct {
	DataRowSorter
	Key string
}

func (this ByKey) Less(i, j int) bool {
	if v1, ok := this.DataRowSorter[i].GetValue(this.Key); ok {
		if v2, ok := this.DataRowSorter[j].GetValue(this.Key); ok {
			return compute.AsFloat64(v1) < compute.AsFloat64(v2)
		}
	}
	return true
}
