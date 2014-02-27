package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type RPCConfig struct {
	DataBases []*DataBase `xml:"DataBases>DataBase"` //多个数据源
	Sqls      []*Sql      `xml:"Sqls>Sql"`           //Sql配置，根据id决定不同的功能
	RPC       *Service    //对外的RPC服务
}

//数据库配置
type DataBase struct {
	Id        string `xml:"id,attr"`
	DSN       string
	Host      string
	Port      int
	SID       string
	User      string
	Password  string
	IsPrimary bool //是否是主库
}

//Sql结构
type Sql struct {
	Id        string   //
	SqlString string   //
	Params    []string `xml:"Params>Param"` //param name
	From      string   //dbid from  database
}

//从文件中读取出Config
func ReadRPC(file string) (*RPCConfig, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result RPCConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//将Config保存为文件
func (this *RPCConfig) Write(file string) error {
	bt, err := xml.Marshal(this)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(file, bt, os.ModePerm); err != nil {
		return err
	}
	return nil
}
