package cache

//import (
//	"monitor/common"
//	"monitor/config"
//	"monitor/input/testinput"
//	"testing"
//)

//func Test_Add(t *testing.T) {
//	mc := NewTimeCache("test", "CID")
//	c, err := config.Read("../config.xml")
//	if err != nil {
//		t.Fatal(err.Error())
//	}
//	Init(c, nil)
//	input := testinput.New(c.Inputs.TestInputs[0])
//	common.IsRun = true
//	input.Start()
//	for rows := range input.Output {
//		mc.Add(rows)
//	}
//}
