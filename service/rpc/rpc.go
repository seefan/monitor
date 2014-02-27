package rpc

import (
	"fmt"
	"monitor/log"
	//	"net/http"
	"monitor/common"
	"net/rpc"
)

var (
	ServerIP string
	Port     int
)

//执行特定id的sql，并返回数据
func BySql(id string) (map[string]int, [][]interface{}, error) {
	if client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", ServerIP, Port)); err != nil {
		var reply = new(common.RPCResponseSql)
		args := common.RPCRequestSql{Id: id}
		err = client.Call("Oracle.BySql", args, &reply)
		if err != nil {
			log.Error("bysql error:", err)
		}
		return reply.Cols, reply.Rows, err
	} else {
		log.Errorf("rcp %s:%d error", ServerIP, Port, err.Error())
	}
	return nil, nil, fmt.Errorf("sql %s not exists", id)
}

//执行带参数的特定id的sql，并返回数据
func BySqlParamName(id string, p map[string]interface{}) (map[string]int, [][]interface{}, error) {
	if client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", ServerIP, Port)); err != nil {
		var reply = new(common.RPCResponseSql)
		args := common.RPCRequestSql{Id: id, Param: p}
		err = client.Call("Oracle.BySqlParamName", args, &reply)
		if err != nil {
			log.Error("bysql error:", err)
		}
		return reply.Cols, reply.Rows, err
	} else {
		log.Errorf("rcp %s:%d error", ServerIP, Port, err.Error())
	}
	return nil, nil, fmt.Errorf("sql %s not exists", id)
}

//执行带参数的特定id的sql
func ExecSqlParamName(id string, p map[string]interface{}) error {
	if client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", ServerIP, Port)); err != nil {
		var reply = new(common.RPCResponseSql)
		args := common.RPCRequestSql{Id: id, Param: p}
		return client.Call("Oracle.ExecuteByParamName", args, &reply)
	} else {
		return fmt.Errorf("rcp %s:%d error", ServerIP, Port, err.Error())
	}
}

//批量执行带参数的特定id的sql
func BatchExecSqlParamName(id string, p []map[string]interface{}) error {
	if client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", ServerIP, Port)); err != nil {
		var reply = new(common.RPCResponseSql)
		args := common.RPCRequestSql{Id: id, Params: p}
		return client.Call("Oracle.BatchExecuteByParamName", args, &reply)
	} else {
		return fmt.Errorf("rcp %s:%d error", ServerIP, Port, err.Error())
	}
}
