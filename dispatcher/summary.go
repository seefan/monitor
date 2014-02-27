package dispatcher

import (
	"monitor/cache"
	"monitor/common"
	"monitor/config"
	"monitor/log"
	"sync"
)

//运算结点，每类数据在一个节点运算，这样需要进行读写锁。
type Summary struct {
	Config             *config.Summary       //配置
	Output             chan *common.DataRows //数据输出流，根据配置由Summary创建
	dict               map[string][]string   //父子关系的逆向字典
	Id                 string
	NeedReloadRelation bool //是否需要重新加载父子关系
	lock               sync.Mutex
}

//创建一个新的节点
func NewSummary(cfg *config.Summary) *Summary {
	log.Infof("node %s is created...", cfg.Id)
	n := &Summary{Config: cfg}
	n.Id = cfg.Id
	n.Output = make(chan *common.DataRows, 100)
	n.init()
	return n
}

func (this *Summary) init() {
	//处理字典
	this.dict = make(map[string][]string)
	this.ReloadRelation()
}

/*
输入数据(小区)根据逆向字典生成输出(场景)数据(也可以使用协程)，每输出(场景)一个协程。(优先，效率高，云计算)
节点间流转数据结构:包括sum，max，min，count，key等的数组为一条记录。按顺序对应字段名。
节点只处理单数据源，无需字段描述，输出时带上数据结构或标明数据结构。
*/

//开始一个新的处理过程
func (this *Summary) Write(row *common.DataRows) {
	if common.IsRun {
		if this.NeedReloadRelation {
			this.ReloadRelation()
			this.NeedReloadRelation = false
		}
		//	log.Infof("get row count", len(row.Rows))
		go this.RunIt(row)
	}
}

//刷新节点关系树
func (this *Summary) ReloadRelation() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if rs, ok := cache.Get("System:Relation:" + this.Id); ok {
		if vs, ok := rs.([]*config.RelationChild); ok {
			this.dict = make(map[string][]string)
			for _, r := range vs { //遍历所有关系表
				for _, key := range r.Childs { //遍历所有子key
					if v, ok := this.dict[key]; ok { //如果字典中存在就增加
						this.dict[key] = append(v, r.Id)
					} else {
						this.dict[key] = []string{r.Id}
					}
				}
			}
		}
	}
}

//close node
func (this *Summary) Close() {
	close(this.Output)
}
func (this *Summary) RunIt(row *common.DataRows) {
	//log.Infoln("收到一批数据", row.Time)
	//将每个输入对应到特定的父级数据
	rmap := make(map[string]*common.DataRow)
	cmap := row.CloneMap()

	for _, v := range row.Rows {
		//find key
		key := row.GetKey(this.Config.Relation.ChildKey, v)
		if pkeys, ok := this.dict[key]; ok {
			for _, pkey := range pkeys {
				if r, ok := rmap[pkey]; ok {
					r.Merge(v)
				} else {
					rrow := new(common.DataRow)
					rrow.RowMap = cmap
					rrow.Time = row.Time
					rrow.Row = v
					rmap[pkey] = rrow
				}
			}
		}
	}
	//合并输出
	result := new(common.DataRows)
	result.Time = row.Time
	result.RowMap = row.CloneMap()
	result.RowMap.Key[this.Config.Relation.PrimaryKey] = cmap.ColumnCount
	result.RowMap.ColumnCount = cmap.ColumnCount + 1
	for k, v := range rmap {
		v.Row = append(v.Row, k)
		result.Rows = append(result.Rows, v.Row)
	}
	if common.IsRun {
		log.Info(result)
		this.Output <- result
	}
}
