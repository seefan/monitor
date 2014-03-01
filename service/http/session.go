package http

import (
	"monitor/log"
	"sync"
	"time"
)

var (
	sessions = make(map[string]*Session) //存贮在线用户的uid和key
	lock     sync.RWMutex
)

const Expired = 60 * 30

//init
func init() {
	go func() {
		timer := time.Tick(time.Minute)
		for t := range timer {
			lock.Lock()
			for k, s := range sessions {
				if time.Since(s.LastTime) > Expired {
					delete(sessions, k)
				}
			}
			lock.Unlock()
			log.Infof("clear session session at %v,used %v", t, time.Since(t))
		}
	}()
}

//web专用session
type Session struct {
	UID      string
	Key      string
	LastTime time.Time
}

//设置值
func SetSession(uid string, key string) {
	s := &Session{UID: uid, Key: key, LastTime: time.Now().Add(time.Duration(Expired))}
	lock.Lock()
	defer lock.Unlock()
	sessions[uid] = s
}

//取取，如果过期，串也为空
func GetSession(id string) (key string, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	if s, ok := sessions[id]; ok && time.Since(s.LastTime).Seconds() < Expired {
		key = s.Key
		ok = true
		s.LastTime = time.Now().Add(time.Duration(Expired))
	}
	return
}

//更新Session的过期时间
func FreshSession(id string) {
	lock.Lock()
	defer lock.Unlock()
	if s, ok := sessions[id]; ok {
		s.LastTime = time.Now().Add(time.Duration(Expired))
	}
}
