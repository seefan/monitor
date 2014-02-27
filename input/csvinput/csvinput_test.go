package csvinput

import (
	"monitor/common"
	"monitor/config"
	"testing"
)

func Test_TestInput(t *testing.T) {
	c, err := config.Read("../../config.xml")
	if err != nil {
		t.Error(err.Error())
	}

	for _, cfg := range c.Inputs.CsvInputs {
		test := New(cfg)
		common.IsRun = true
		test.Config.FileName = "../../test.csv"
		test.Start()
		for o := range test.Output {
			println(o.Time)
		}
	}
}
