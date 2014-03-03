package protocol

import (
	"fmt"
	"monitor/common"
	"monitor/log"
	"monitor/service/online"
	"monitor/service/rpc"
	"time"
)

//pid process
type PidProcess func(cmd *common.RequestData, c *online.Client)

var (
	pid map[int]PidProcess //all pid
	now time.Time          //start time
)

//init pid
func init() {
	pid = map[int]PidProcess{ //以0为区分，0以下的id都不需要进行登陆
		0:  pid0,  //运行状态测试
		-1: pid_1, //登陆
		-2: pid_2, //退出
		-3: pid_3, //用临时key重新登陆，用户登陆后会临时产生
		2:  pid2,  //数据推送注册
		3:  pid3,  //数据库的读操作
		4:  pid4,  //数据库的写操作
		5:  pid5,  //更新节点内的父子关系
	}
	now = time.Now()
}

//exec pid function
func Exec(cmd *common.RequestData, c *online.Client) {

	//没有登陆就必须登陆系统，小于等于0的协议不需要登陆。
	if cmd.Pid <= 0 {
		pid[cmd.Pid](cmd, c)
	} else if !c.IsLogin {
		cmd := common.ResponseData{1, 0, "need login", nil}
		c.Send(&cmd)
	} else {
		if p, ok := pid[cmd.Pid]; ok {
			p(cmd, c)
		} else {
			pidok(cmd, c)
		}
	}
}

//数据库read操作
func pid3(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "Id", "Param") {
		return
	}
	id := cmd.GetString("Id")
	m := make(map[string]interface{})
	col, rows, err := rpc.BySqlParamName(id, m)
	if err != nil {
		log.Error(err.Error())
		writePrror("db error", cmd, c)
		return
	}
	dp := new(common.OutputParam)
	dp.RowMap = col
	dp.Rows = rows
	dp.Time = time.Now()
	writepid(c, cmd, 0, dp)
}

//数据库write操作
func pid4(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "Id", "Param") {
		return
	}
	id := cmd.GetString("Id")
	m := make(map[string]interface{})
	err := rpc.ExecSqlParamName(id, m)
	if err != nil {
		log.Error(err.Error())
		writePrror("db error", cmd, c)
		return
	}
	writepid(c, cmd, 0, nil)
}

//注册推送数据的请求
func pid2(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "DataKey", "DataRange", "Id", "LimitId", "OrderBy", "OrderKey", "OutputKey", "TimeRange") {
		return
	}
	dp := new(common.DataParam)
	dp.DataKey = cmd.GetString("DataKey")
	dp.DataRange = cmd.GetInt("DataRange")
	dp.Id = cmd.GetString("Id")
	dp.LimitId = cmd.GetStringArray("LimitId")
	dp.OrderBy = cmd.GetInt("OrderBy")
	dp.OrderKey = cmd.GetString("OrderKey")
	dp.OutputKey = cmd.GetStringArray("OutputKey")
	dp.FillKey = cmd.GetString("FillKey")
	dp.TimeRange = cmd.GetInt("TimeRange")
	c.AddRequest(dp)
	//同时直接返回一次数据，否则界面会有空白期
	tm := &common.TimeMessage{Time: online.GetCurrentTime()}
	c.Processing(tm)
}

//更新节点内的父子关系，有3类变化，新增、删除
func pid5(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "DataKey", "Insert", "Delete", "Id") {
		return
	}
}

//尚未支持
func pidok(cmd *common.RequestData, c *online.Client) {
	writepid(c, cmd, -1, "not surport")
}

// test service
func pid0(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "Test") {
		return
	}

	test := cmd.GetInt("Test")
	var result interface{}
	switch test {
	case 0:
		result = true
	case 1:
		result = time.Now().Format(common.TimeFormat)
	case 2:
		result = time.Since(now).String()
	case 3:
		result = now.Format(common.TimeFormat)
	default:
		result = fmt.Sprintf("not surport %d", test)
	}
	writepid(c, cmd, 0, result)
}

//check param  exists
func checkParamError(cmd *common.RequestData, c *online.Client, keys ...string) bool {
	for _, k := range keys {
		if _, ok := cmd.Param[k]; !ok {
			writePrror(fmt.Sprintf("param '%s' not exists", k), cmd, c)
			return true
		}
	}
	return false
}

//return error status=-1
func writePrror(msg string, cmd *common.RequestData, c *online.Client) {
	writepid(c, cmd, -1, msg)
}

//key login in
func pid_3(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "UID", "Key") {
		return
	}

	login := cmd.GetString("UID")
	pwd := cmd.GetString("Key")

	loginSuccess := false
	//log.Info("pid -3 uid=key", login, pwd)
	//check session
	if name, key, ok := online.GetSession(login); ok && key == pwd {
		c.Name = name
		c.Key = key
		loginSuccess = true
		//log.Info("get session success", key)
	}

	//check
	re := make(map[string]interface{})
	if loginSuccess {
		c.IsLogin = true
		re["LoginState"] = 0
		re["Message"] = "login success"
		re["UID"] = c.UUID
		re["Key"] = c.Key
	} else {
		re["LoginState"] = -1
		re["Message"] = "login failed"
	}
	online.Set(c)
	writepid(c, cmd, 0, re)
}

// login in
func pid_1(cmd *common.RequestData, c *online.Client) {
	if checkParamError(cmd, c, "Login", "Password") {
		return
	}

	login := cmd.GetString("Login")
	pwd := cmd.GetString("Password")
	c.Name = login
	log.Infof("user %s login", c.Name, common.IsDebug)
	loginSuccess := false
	if common.IsDebug { //调试时直接登陆，不验证密码
		loginSuccess = true
	} else {
		param := make(map[string]interface{})
		param["login"] = login
		_, epwd, err := rpc.BySqlParamName("login", param)
		if err != nil {
			log.Error(err.Error())
		}
		if err == nil && len(epwd) > 0 {
			loginSuccess = epwd[0][0] == pwd //需要对密码进行加密再进行比较
		}
	}

	//check
	re := make(map[string]interface{})
	if loginSuccess {
		c.IsLogin = true
		c.Key = common.HashString(c.UUID + ":" + c.Name)
		online.SetSession(c.UUID, c.Name, c.Key)
		re["LoginState"] = 0
		re["Message"] = "login success"
		re["UID"] = c.UUID
		re["Key"] = c.Key
		online.Set(c)
	} else {
		re["LoginState"] = -1
		re["Message"] = "login failed"
	}
	writepid(c, cmd, 0, re)
}

// //senc data to online.Client
// func pid2(cmd *common.RequestData, c *online.Client) {
// 	log.Println("pid 2")
// 	//if param, ok := cmd.Param.(string); ok {
// 	if start_time, err := time.Parse(common.TimeFormat, cmd.Param["time"]); err == nil {
// 		log.Println(start_time)
// 		out := online.ResponseData{1, 0, nil, cmd.Param}
// 		//get from redis and send to all online.Client
// 		this.SendToAll(out)
// 	}
// }

//login out
func pid_2(cmd *common.RequestData, c *online.Client) {
	//if param, ok := cmd.Param.(string); ok {
	log.Debug("login out ", c.Name)
	writepid(c, cmd, 0, nil)
	c.IsRun = false
	online.RemoveSession(c.UUID)
	//}
}

//write pid to responseout ResponseData
func writepid(c *online.Client, cmd *common.RequestData, status int, result interface{}) {
	out := common.ResponseData{cmd.Pid, status, result, cmd.Param}
	c.Send(&out)
}
