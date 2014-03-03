package cache

import (
	"bytes"
	"container/ring"
	"encoding/gob"
	"fmt"
	"github.com/alphazero/Go-Redis"
	"monitor/common"
	"monitor/config"
	"monitor/log"
	"sync"
	"time"
)

var (
	pool        = make(map[string]chan *RedisClient)
	errRedis    = make(map[string]*config.Redis) //重连机制，定时重试创建全新的chan，成功后从这里删除
	keyToClient = make(map[byte]string)
	ids         *ring.Ring
	redislock   sync.RWMutex
)

func RedisStatus() {
	fmt.Println("pool size:", len(pool))
}

type RedisClient struct {
	Redis redis.Client
	Id    string
}

//id是否可用
type IdAvailable struct {
	Id          string
	IsAvailable bool
}

func (this *RedisClient) Close() {
	//println("close client at %s", this.Id)
	pool[this.Id] <- this
}

const (
	PoolSize = 4
)

//初始化redis连接池
func initRedis(confs []*config.Redis) {
	ids = ring.New(len(confs))
	for _, conf := range confs {
		if !conf.Enable {
			continue
		}
		id := &IdAvailable{Id: conf.Id}
		ids.Value = id
		spec := redis.DefaultSpec().Db(conf.DB).Password(conf.Password).Host(conf.Host).Port(conf.Port)
		//log.Info("redis init at ", spec)
		c := make(chan *RedisClient, PoolSize)
		for j := 0; j < PoolSize; j++ {
			if client, err := redis.NewSynchClientWithSpec(spec); err != nil {
				goto err
			} else {
				c <- &RedisClient{Id: conf.Id, Redis: client}
			}
		}
		pool[conf.Id] = c
		id.IsAvailable = true
		ids = ids.Next()
	err: //如果创建CLIENT时出错，就抛弃这个台机器
		errRedis[conf.Id] = conf
		ids = ids.Next()
	}
}

//取一个client的id，如果没有映射就要新增
func getId(key string) (string, error) {
	bytes := []byte(key)
	if len(bytes) > 0 {
		redislock.Lock()
		defer redislock.Unlock()
		b := bytes[0]
		//log.Info("client b is %s", b)
		if id, ok := keyToClient[b]; ok { //如果有记录有对应关系
			//log.Info("client id is %s", id)
			return id, nil
		} else {
			//首先检查是否
			for tid := ids.Value; tid != nil && tid.(*IdAvailable).IsAvailable; ids = ids.Next() {
				id = tid.(*IdAvailable).Id
				//log.Info("client id is %s", id)
				break
			}
			return id, nil
		}
	}
	return "", fmt.Errorf("error get id with %s", key)
}

//创建一个连接
func GetClient(key string) (*RedisClient, error) {
	if id, err := getId(key); err == nil {
		if c, ok := pool[id]; ok {
			timeout := time.After(time.Second * 10)
			//client := <-c

			select {
			case client := <-c:
				//log.Info("get client id is %s", client.Id)
				return client, nil
			case <-timeout:
				return nil, fmt.Errorf("get connection time out")
			}
			return <-c, nil
		}
	}
	return nil, fmt.Errorf("get client error ")
}

//将数据保存到外部缓存中，同时也支持从外部缓存取数据，在这里进行屏蔽，其它地方不操作外部缓存
//保存到外部存贮
func saveToStore(key string, value *common.DataRow, expired int) {
	//begin := time.Now()
	if len(pool) > 0 {
		if c, err := GetClient(value.Time); err == nil {
			defer c.Close()
			if bs, err := encode(*value); err == nil {
				c.Redis.Hset(key, value.Time, bs)
				c.Redis.Expire(key, int64(expired*60)) //秒
			}
		} else {
			log.Error("get client error,begin recycle", err.Error())
			//回收连接，过段时间再试
		}
	}
	//diff := time.Since(begin)
	//log.Infof("save to store at %v seconds", diff.Seconds())
}

//从外部取数据
func getFromStore(key string, keyTime string) (*common.DataRow, bool) {
	if c, err := GetClient(keyTime); err == nil {
		defer c.Close()
		if bts, err := c.Redis.Hget(key, keyTime); err == nil {
			if r, err := decode(bts); err == nil {
				return r, true
			}
		}
	}
	return nil, false
}

//数据编码
func encode(vs common.DataRow) ([]byte, error) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(vs)
	return buff.Bytes(), err
}

//数据解码
func decode(vs []byte) (*common.DataRow, error) {
	buff := bytes.NewBuffer(vs)
	dec := gob.NewDecoder(buff)
	result := new(common.DataRow)
	err := dec.Decode(result)
	return result, err
}
