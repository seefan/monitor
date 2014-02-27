package csvinput

import (
	"io/ioutil"
	"monitor/common"
	"monitor/config"
	"monitor/input"
	"monitor/log"
	"strconv"
	"strings"
	"time"
)

type CsvInput struct {
	input.InputFace
	Config *config.CsvInput
	TIME   time.Time
}

//start oracle face
func New(cfg *config.CsvInput) (db *CsvInput) {
	db = new(CsvInput)
	db.Config = cfg
	db.TIME = time.Now().Truncate(time.Minute)
	db.Output = make(chan *common.DataRows, 100)
	db.InitColumn(cfg.OutputColumn, cfg.OutputKey)
	return
}

func (this *CsvInput) Start() {
	go this.Run()
}
func (this *CsvInput) GetOutput() chan *common.DataRows {
	return this.Output
}
func (this *CsvInput) GetConfig() *config.Input {
	return &this.Config.Input
}

//run service
func (this *CsvInput) Run() {
	log.Infoln("CsvInput is starting")
	log.Infof("open file %s", this.Config.FileName)

	g := float64(this.Config.Period)
	if bts, err := ioutil.ReadFile(this.Config.FileName); err == nil {
		strs := strings.Split(string(bts), "\n")
		for common.IsRun {
			beginget := time.Now()
			mod := this.TIME.Minute() % 3
			//			log.Infof("w count is %d", mod+1)
			for i := 0; i <= mod; i++ {
				row := new(common.DataRows)
				row.Time = time.Now().Format(this.Config.TimeFormat)
				row.RowMap = this.RowMap
				for _, str := range strs {
					if this.Config.OutputColumn == nil {
						log.Infoln("output column is nil")
						break
					}
					cols := strings.Split(str, ",")
					var r []interface{}
					for _, c := range cols {
						if d, err := strconv.ParseFloat(c, 10); err == nil {
							r = append(r, d)
						} else {
							r = append(r, 0)
						}
					}
					row.Rows = append(row.Rows, r)
				}
				this.Output <- row
				//log.Infof("push row is %v", row.Time)
			}
			//计算下一次运行的时间

			this.TIME = this.TIME.Add(time.Minute * time.Duration(g))
			log.Infof("next time is %v", this.TIME)
			sec := time.Since(beginget).Seconds()

			//计算下一次运行的时间
			if sec < g*60 { //如果一个处理周期内处理完，就延时，没有处理完就立即执行,如果正式运行期间还需要进行时间纠正
				log.Infof("sleep seconds %f", (g*60)-sec)
				time.Sleep(time.Second * (time.Duration(g*60 - sec))) //sleep one minute-sec
			}
		}

	} else {
		log.Error(err.Error())
	}
	close(this.Output)
}
