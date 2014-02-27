package common

import (
	"fmt"
)

type RPCResponseSql struct {
	Cols map[string]int
	Rows [][]interface{}
}
type RPCRequestSql struct {
	Id     string
	Param  map[string]interface{}
	Params []map[string]interface{}
}

func (this *RPCRequestSql) ToString() string {
	return fmt.Sprintf("Sql Id:%s\tParam:%v", this.Id, this.Param)
}
