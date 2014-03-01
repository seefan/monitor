package main

import (
	//	"io/ioutil"
	"monitor/cache"
	"monitor/common"
	"monitor/config"
	"monitor/dispatcher"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(common.MaxCPU)
	println("cpu num is ", common.MaxCPU)
	cfg, err := config.Read("config.xml")
	if err != nil {
		panic(err.Error())
	}
	rcfg, err := config.ReadRelation("relation.xml")
	if err != nil {
		panic(err.Error())
	}
	cache.Init(cfg, rcfg)
	LoadBaseData()
	dispatcher.Start(cfg)
	time.Sleep(time.Hour)
	//if str, err := ioutil.ReadFile("config/config.xml"); err == nil {
	//	ioutil.WriteFile("config.xml", str, os.ModePerm)
	//} else {
	//	panic(err.Error())
	//}//

	//test.Test_Dispatcher()
	//test.Test_MutiCache()
	//test.Test_Push()
	//test.Test_WebSocket()
	a, err := time.Parse(common.TimeFormat, "2014-02-12 12:59")
	b, err := time.Parse(common.TimeFormat, "2014-02-12 13:00")
	println(a.Before(b), err)
	println(a.After(b))
	println(a.Sub(b))
	println(a.String())
	println(b.String())
}
func LoadBaseData() {
	cache.Set("Relation:1", []string{"1", "2", "3", "4", "5"})
	cache.Set("Relation:11", []string{"1", "2", "3", "4", "5"})
	cache.Set("Relation:2", []string{"1", "2", "3", "4", "5"})
}
