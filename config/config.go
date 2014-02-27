package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

//配置
type Config struct {
	Redises   []*Redis   `xml:"Redises>Redis"` //多个缓存
	Inputs    *DataFace  //多个输入
	Summarys  []*Summary `xml:"Summarys>Summary"` //多个计算单元
	Http      *Service   //对外的HTTP服务
	Tcp       *Service   //对外的TCP服务
	RPC       *Service   //外部的RPC服务
	PushDelay float64    //发送数据的时间间隔，如果设置为0,就进行实时推送，单位为秒
	IsDebug   bool       //是否为调试状态，[调试时直接登陆，不验证密码；]
}

//将Config转成xml格式
func (this *Config) String() string {
	if bt, err := xml.Marshal(this); err == nil {
		return string(bt)
	} else {
		return err.Error()
	}
}

//将Config直接转成字节数组
func (this *Config) Bytes() []byte {
	if bt, err := xml.Marshal(this); err == nil {
		return bt
	}
	return nil
}

//从文件中读取出Config
func Read(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result Config
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//将Config保存为文件
func (this *Config) Write(file string) error {
	bt, err := xml.Marshal(this)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(file, bt, os.ModePerm); err != nil {
		return err
	}
	return nil
}

//支持的数据输入接口
type DataFace struct {
	TcpInputs   []*TcpInput
	TestInputs  []*TestInput
	CsvInputs   []*CsvInput
	RedisInputs []*RedisInput
}

//一个数据库的输入项
type SqlInput struct {
	Input
	SqlId string //sql语句，用于查询数据库，注意与输出列对应
}

//一个测试用的输入项，用于测试系统运行状况
type TestInput struct {
	Input
	TimeKey string //时间字段的名称，默认为START_TIME，用于推动数据流动
}

//csv文件输入，用于测试数据准确性
type CsvInput struct {
	Input
	FileName string //文件名
}

//从Redis实时接收数据
type RedisInput struct {
	Input
	From  string //redis的id
	SubId string //订阅的id
}

//基础输入参数结构
type Input struct {
	Id           string      `xml:"id,attr"`
	OutputKey    []string    `xml:"OutputKeys>Key"`       //输出的关键字段key
	OutputColumn *Expression `xml:"OutputColumns>Column"` //输出列
	Enable       bool        //是否启用
	EnableStore  bool        //是否进行存贮
	PrimaryKey   *Key        //主键字段名
	TimeFormat   string      //时间的格式串
	Period       int         //粒度，以分钟为单位，重复输入的时间间隔。主要用于处理时间字段。
}

//tcp的输入参数
type TcpInput struct {
	Input
	InputKey []string //输入字段的列名，按位置对应
}

//通用key，value
type Key struct {
	Type  string `xml:"type,attr"` //类型，包括int float time string
	Value string
}

//汇总单元，对数据按指定规则进行汇聚
type Summary struct {
	Id          string `xml:"id,attr"`
	From        string
	Outputs     []*Key    `xml:"Outputs>Output"` //输出的key
	Relation    *Relation //该计算单元所有需要计算的，类似bsc->cell
	TimeName    string    //时间字段的名称，默认为START_TIME，用于推动数据流动
	EnableStore bool      //是否进行存贮
}

//对应的父子关系
type Relation struct {
	PrimaryKey string //主键名
	ChildKey   string //对应子级的字段名，与PrimaryKey联合使用
}

//指标计算方式
type Expression struct {
	Sum []string
	Max []string
	Min []string
}

//Redis配置
type Redis struct {
	Id       string `xml:"id,attr"`
	Host     string
	Port     int
	Password string
	DB       int
}

//Tcp服务
type Service struct {
	Host   string //ip地址
	Port   int    //端口
	Enable bool   //是否启用
}
