package config

import (
	"testing"
)

func Test_config(t *testing.T) {
	cfg := new(Config)
	exp := new(Expression)
	exp.Sum = []string{"N2", "N3", "N6"}
	baseInput := Input{
		Id:           "ACELL",
		Enable:       true,
		OutputKey:    []string{"CID", "START_TIME"},
		OutputColumn: exp,
		TimeFormat:   "2006-01-02 15:04",
		PrimaryKey:   &Key{Type: "number", Value: "CID"},
		EnableStore:  true,
		Period:       1,
	}
	cfg.Inputs = new(DataFace)

	rs := []*Redis{&Redis{Host: "192.168.8.113", Port: 6379, DB: 12, Id: "redis113"}}
	cfg.Redises = rs

	ga := new(Summary)
	ga.Id = "BSC"
	ga.From = "ACELL"
	ga.TimeName = "START_TIME"
	ga.EnableStore = true
	ga.Relation = &Relation{ChildKey: "CID", PrimaryKey: "BSC"}
	cfg.Summarys = []*Summary{ga}
	test := &TestInput{
		Input:   baseInput,
		TimeKey: "START_TIME",
	}
	test.Enable = false
	csv := &CsvInput{
		Input:    baseInput,
		FileName: "test.csv",
	}
	csv.Enable = true
	csv.EnableStore = true
	cfg.Inputs.TestInputs = []*TestInput{test}
	cfg.Inputs.CsvInputs = []*CsvInput{csv}
	cfg.Http = &Service{Host: "0.0.0.0", Port: 6780, Enable: true}
	cfg.Tcp = &Service{Host: "0.0.0.0", Port: 6781, Enable: true}
	cfg.RPC = &Service{Host: "0.0.0.0", Port: 6782, Enable: true}
	cfg.PushDelay = 30
	cfg.IsDebug = true
	if err := cfg.Write("../config.xml"); err != nil {
		t.Error(err.Error())
	}
	if cfg, err := Read("../config.xml"); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(cfg)
	}
}
func Test_rpc(t *testing.T) {
	cfg := new(RPCConfig)
	cfg.DataBases = []*DataBase{&DataBase{Id: "db113", Host: "192.168.8.113", Port: 1521, User: "dsam", Password: "dsam", SID: "orcl", IsPrimary: true}}
	cfg.Sqls = []*Sql{
		&Sql{Id: "acell", SqlString: "SELECT sum(n2) n2,sum(n3) n3,sum(n6) n6,cid, to_char( start_time,'yyyy-mm-dd hh24:mi') start_time from DS_PRD_A_CELL_M t where start_time>=to_date(':start_time','yyyy-mm-dd hh24:mi:ss') and start_time<to_date(':end_time','yyyy-mm-dd hh24:mi:ss') group by cid,to_char( start_time,'yyyy-mm-dd hh24:mi')", Params: []string{"start_time", "end_time"}, From: "db113"},
		&Sql{Id: "login", SqlString: "SELECT USER_PWD FROM CFG_USER_G WHERE USER_NAME = :login", Params: []string{"login"}},
		&Sql{Id: "all_node", SqlString: "SELECT SCENEID, NEID FROM CFG_SCENE_NE", Params: []string{}},
		&Sql{Id: "insert_node", SqlString: "insert into CFG_SCENE_NE (SCENEID, NEID) values(:parentid,:childid) ", Params: []string{"parentid", "childid"}},
		&Sql{Id: "delete_node", SqlString: "delete from CFG_SCENE_NE where  SCENEID=:parentid and  NEID=:childid) ", Params: []string{"parentid", "childid"}},
	}
	cfg.RPC = &Service{Host: "0.0.0.0", Port: 6782, Enable: true}
	if err := cfg.Write("../rpc.xml"); err != nil {
		t.Error(err.Error())
	}
	if cfg, err := Read("../rpc.xml"); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(cfg)
	}
}
func Test_relation(t *testing.T) {
	cfg := new(RelationConfig)
	d1 := &RelationChild{Id: "1", Childs: []string{"1", "10", "100"}}
	d2 := &RelationChild{Id: "1", Childs: []string{"1", "3", "100"}}
	cfg.Relations = []*RelationParent{&RelationParent{Id: "BSC", Childs: []*RelationChild{d1, d2}}}
	cfg.Write("../relation.xml")
	if _, err := ReadRelation("../relation.xml"); err != nil {
		t.Error(err.Error())
	}
}
