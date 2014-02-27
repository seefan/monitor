package testinput

import (
	"math/rand"
	"monitor/common"
	"monitor/config"
	"monitor/input"
	"monitor/log"
	"time"
)

type TestInput struct {
	input.InputFace
	Config *config.TestInput
	TIME   time.Time
}

//start oracle face
func New(cfg *config.TestInput) (db *TestInput) {
	db = new(TestInput)
	db.Config = cfg
	db.Output = make(chan *common.DataRows, 100)
	db.TIME = time.Now().Truncate(time.Minute)
	db.InitColumn(cfg.OutputColumn, cfg.OutputKey)
	return
}
func (this *TestInput) Start() {
	go this.Run()
}
func (this *TestInput) GetOutput() chan *common.DataRows {
	return this.Output
}
func (this *TestInput) GetConfig() *config.Input {
	return &this.Config.Input
}

//run service
func (this *TestInput) Run() {
	log.Infoln("TestInput is starting")
	g := float64(this.Config.Period)
	for common.IsRun {
		beginget := time.Now()
		for i := 0; i < 60; i++ {
			row := new(common.DataRows)
			row.RowMap = this.RowMap
			if this.Config.OutputColumn == nil {
				log.Infoln("output column is nil")
				break
			}
			rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
			for j := 0; j < 100000; j++ {
				var r []interface{}
				if this.Config.OutputColumn.Sum != nil {
					for k := 0; k < len(this.Config.OutputColumn.Sum); k++ {
						r = append(r, rnd.Float32())
					}
				}
				if this.Config.OutputColumn.Max != nil {
					for k := 0; k < len(this.Config.OutputColumn.Max); k++ {
						r = append(r, rnd.Float32())
					}
				}
				if this.Config.OutputColumn.Min != nil {
					for k := 0; k < len(this.Config.OutputColumn.Min); k++ {
						r = append(r, rnd.Float32())
					}
				}

				if this.Config.OutputKey != nil {
					for _, k := range this.Config.OutputKey {
						var v interface{}
						if k == this.Config.TimeKey {
							v = time.Now()
						} else {
							v = i * j
						}
						r = append(r, v)
					}
				}
				row.Rows = append(row.Rows, r)
			}

			row.Time = this.TIME.Format(this.Config.TimeFormat)

			//log.Infof("testinput 正在输出数据...%v", row.Time)
			this.Output <- row

		}
		//计算下一次运行的时间

		this.TIME = this.TIME.Add(time.Minute * time.Duration(g))
		log.Infof("next time is %v", this.TIME)
		//计算时间间隔
		sec := time.Since(beginget).Seconds()

		//计算下一次运行的时间
		if sec < g*60 { //如果一个处理周期内处理完，就延时，没有处理完就立即执行,如果正式运行期间还需要进行时间纠正
			log.Infof("sleep seconds %f", (g*60)-sec)
			time.Sleep(time.Second * (time.Duration(g*60 - sec))) //sleep one minute-sec
		}

	}
	close(this.Output)
}
