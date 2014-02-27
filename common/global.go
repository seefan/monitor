package common

import (
	"runtime"
	"time"
)

var (
	MaxCPU      int           //最大可用的CPU数
	StoreMinute time.Duration = 2
	StoreCount  int           = 20
	IsRun       bool          //系统是否正在运行
	IsDebug     bool          //是否为调试状态
)

const (
	TimeFormat     string = "2006-01-02 15:04"
	MiniTimeFormat string = "02:15:04"
	NewLine        byte   = '\n'
	BufferSize            = 1024
)

func init() {
	MaxCPU = runtime.NumCPU() - 1
}
