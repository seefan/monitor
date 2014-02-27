package cache

import (
	"fmt"
	"monitor/common"
	"monitor/config"
	"strings"
	"testing"
)

func Test_Redis(t *testing.T) {
	cfg, err := config.Read("../config.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	Init(cfg, nil)
	if c, err := GetClient("d"); err != nil {
		t.Error(err.Error())
	} else {
		RedisStatus()
		if result, err := c.Redis.AllKeys(); err == nil {
			t.Log("allkey", strings.Join(result, ","))
		} else {
			t.Error(err.Error())
		}
		if err := c.Redis.Set("test", []byte("dfdsiedskioe")); err != nil {
			t.Error(err.Error())
		}
		if ss, err := c.Redis.Get("test"); err != nil {
			t.Error(err.Error())
		} else {
			t.Log("test=", string(ss))
		}
		//for i := 0; i < 100000; i++ {
		//	c.Redis.Hset("tt", fmt.Sprintf("%s", i), []byte("test"))
		//}
		//for i := 0; i < 100000; i++ {
		//	c.Redis.Set(fmt.Sprintf("teaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadsafsdafdsfdsafasdfdsafdsafsadfasfasfasfasfaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaast%s", i), []byte("test"))
		//}
		c.Close()

		for i := 0; i < 100000; i++ {
			if cc, err := GetClient(fmt.Sprintf("%s", i)); err != nil {
				t.Error(err.Error())
			} else {
				cc.Redis.Set(fmt.Sprintf("%s:%d", "teaaaaaaaaaaaaaaaaaaaaaaaaaaaaaadsafsdafdsfdsafasdfdsafdsafsadfasfasfasfasfaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaast", i), []byte("test"))
				cc.Close()
			}

		}

	}

}
func Test_Code(t *testing.T) {
	row := &common.DataRow{}
	row.Row = append(row.Row, 1)
	row.RowMap = new(common.DataRowMap)
	row.RowMap.Key = make(map[string]int)
	row.RowMap.Key["id"] = 0
	v, _ := encode(*row)
	vr, _ := decode(v)
	println("vr=", vr.RowMap.Key)
}
