package dispatcher

import (
	"monitor/cache"
	"monitor/common"
	"monitor/config"
	"monitor/input"
	"monitor/input/csvinput"
	"monitor/input/testinput"
	"monitor/log"
	"monitor/service"
)

var (
	inputs   = make(map[string]input.InputStartFace)
	summarys = make(map[string]*Summary)
	link     = make(map[string][]string)
)

//启动处理引擎
func Start(conf *config.Config) {
	log.Infof("dispatcher is starting")
	common.IsRun = true
	common.IsDebug = conf.IsDebug
	createDataPeriod(conf)
	service.Start(conf)
	createInput(conf)
	createSummary(conf)
	//cron.Start(conf)
	Run()
}

//创建数据类型与时间粒度的关系
func createDataPeriod(conf *config.Config) {
	for _, c := range conf.Inputs.CsvInputs {
		createPeriodFromInput(&c.Input)
	}
	for _, c := range conf.Inputs.RedisInputs {
		createPeriodFromInput(&c.Input)
	}
	for _, c := range conf.Inputs.TcpInputs {
		createPeriodFromInput(&c.Input)
	}
	for _, c := range conf.Inputs.TestInputs {
		createPeriodFromInput(&c.Input)
	}
}
func createPeriodFromInput(conf *config.Input) {
	cache.Set(cache.FormatKey("System", "Period", conf.Id), conf.Period)
}

//关闭处理引擎
func Close() {
	common.IsRun = false
	for k, _ := range inputs {
		delete(inputs, k)
	}
	for k, _ := range summarys {
		delete(summarys, k)
	}
	cache.CloseCache()
}

//按配置把输入和输出连接起来
func Run() {
	//接收时间通知事件，将时间变化通知给服务
	cache.TimeChanged = service.MessageTimeChange

	for k, v := range inputs {
		go runInput(k, v)
	}
	for k, v := range summarys {
		go runSummary(k, v)
	}

}

//处理数据输入流
func runInput(inputName string, in input.InputStartFace) {
begin:
	out := in.GetOutput()
	for r := range out { //收到一条输出数据，两种选择，一是传给需要的节点，二是保存到缓存
		go func(r *common.DataRows) {
			if ns, ok := link[inputName]; ok { //如果有节点需要本节点的输出
				for _, fname := range ns { //给所有节点
					summarys[fname].Write(r) //写数据
				}
			}
			//如果需要保存到CACHE，这里处理
			//log.Infof("input is %v", in.GetConfig().EnableStore, inputName, r.Time)
			cache.AddRowsToCache(inputName, r)
		}(r)
	}
	if common.IsRun { //如果程序还没有结束，就重新启动INPUT
		in.Start()
		goto begin
	}
}

//处理节点去运算
func runSummary(inputName string, node *Summary) {
	for rows := range node.Output { //接收数据
		if ns, ok := link[inputName]; ok { //如果有节点需要本节点的输出
			for _, fname := range ns { //给所有节点
				summarys[fname].Write(rows) //写数据
			}
		}
		//如果需要保存到CACHE，这里处理
		//	log.Infof("node is %v", node.Config.EnableStore, inputName, rows.Time)

		cache.AddRowsToCache(inputName, rows)

	}
}

//按参数创建所有的输入，并加到一个字典里
func createSummary(conf *config.Config) {
	for _, n := range conf.Summarys {
		summarys[n.Id] = NewSummary(n) //创建一个节点
		makeLink(n.Id)
	}
}

//创建节点间，数据输入流单的映射关系
func makeLink(id string) {
	for _, fn := range summarys { //处理节点之间的引用关系
		if fn.Config.From == id { //找到引用当前节点的节点
			if ns, ok := link[id]; ok {
				ns = append(ns, fn.Id)
				link[id] = ns
			} else {
				link[id] = []string{fn.Id}
			}
		}
	}
}

//按参数创建所有的输入，并加到一个字典里
func createInput(conf *config.Config) {
	for _, s := range conf.Inputs.TestInputs {
		if !s.Enable {
			continue
		}
		in := testinput.New(s)
		in.Start()

		inputs[s.Id] = in
		makeLink(s.Id)
	}
	for _, s := range conf.Inputs.CsvInputs {
		if !s.Enable {
			continue
		}
		in := csvinput.New(s)
		in.Start()

		inputs[s.Id] = in
		makeLink(s.Id)
	}
}
