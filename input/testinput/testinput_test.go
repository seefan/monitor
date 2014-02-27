package testinput

import (
	"monitor/config"
	"testing"
	"time"
)

func Test_TestInput(t *testing.T) {
	c, err := config.Read("../../config.xml")
	if err != nil {
		t.Error(err.Error())
	}
	for _, cfg := range c.Inputs.TestInputs {
		test := New(cfg)
		go func() {
			for out := range test.Output {
				t.Log(out)
			}
		}()
		test.Start()
		time.Sleep(time.Second)
	}
}
