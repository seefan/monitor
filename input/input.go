package input

import (
	"monitor/common"
	"monitor/config"
	"monitor/log"
)

type InputStartFace interface {
	Start()
	GetConfig() *config.Input
	GetOutput() chan *common.DataRows
}
type InputFace struct {
	Output chan *common.DataRows
	RowMap *common.DataRowMap
}

func (this *InputFace) InitColumn(c *config.Expression, k []string) {
	this.RowMap = new(common.DataRowMap)
	if c == nil {
		log.Infoln("output column is nil")
		return
	}
	i := 0
	this.RowMap.Key = make(map[string]int)
	if c.Sum != nil {
		this.RowMap.Sum = make(map[string]int)
		for _, v := range c.Sum {
			this.RowMap.Sum[v] = i
			this.RowMap.Key[v] = i
			i = i + 1
		}
	}
	if c.Max != nil {
		this.RowMap.Max = make(map[string]int)
		for _, v := range c.Max {
			this.RowMap.Key[v] = i
			this.RowMap.Max[v] = i
			i = i + 1
		}
	}
	if c.Min != nil {
		this.RowMap.Min = make(map[string]int)
		for _, v := range c.Min {
			this.RowMap.Key[v] = i
			this.RowMap.Min[v] = i
			i = i + 1
		}
	}

	if k != nil {
		for _, v := range k {
			this.RowMap.Key[v] = i
			i = i + 1
		}
	}
}
