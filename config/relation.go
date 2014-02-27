package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

//数据关系，用于标记一类数据的父子关系
type RelationChild struct {
	Id     string
	Childs []string `xml:"Childs>Child"`
}

//数据关系配置
type RelationConfig struct {
	Relations []*RelationParent `xml:"Relations>Relation"`
}

//数据关系，用于标记一类数据的父子关系的集合
type RelationParent struct {
	Id     string
	Childs []*RelationChild `xml:"Childs>Child"`
}

//从文件中读取出Config
func ReadRelation(file string) (*RelationConfig, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var result RelationConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

//将Config保存为文件
func (this *RelationConfig) Write(file string) error {
	bt, err := xml.Marshal(this)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(file, bt, os.ModePerm); err != nil {
		return err
	}
	return nil
}
