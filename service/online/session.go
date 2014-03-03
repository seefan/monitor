package online

import (
	"monitor/log"
	"time"
)

var (
	sessions = make(map[string]*Session) //存贮在线用户的uid和key
)

const Expired = 60 * 30

//init
func init() {
	log.Info("Session is started,expired every %d seconds", Expired)
	go func() {
		var i int
		timer := time.Tick(time.Minute)
		for t := range timer {
			i = 0
			lock.Lock()
			for k, s := range sessions {
				if s.IsExpired() {
					delete(sessions, k)
					i = i + 1
				}
			}
			lock.Unlock()
			log.Infof("clear %d sessions at %v,used %v", i, t, time.Since(t))
		}
	}()
}

//web专用session
type Session struct {
	UID      string
	Name     string
	Key      string
	LastTime time.Time
}

//是否过期
func (this *Session) IsExpired() bool {
	return time.Since(this.LastTime).Seconds() > Expired
}

//设置值
func SetSession(uid string, name string, key string) {
	s := &Session{UID: uid, Key: key, LastTime: time.Now()}
	lock.Lock()
	defer lock.Unlock()
	sessions[uid] = s
}

//取取，如果过期，串也为空
func GetSession(id string) (name string, key string, ok bool) {
	lock.RLock()
	defer lock.RUnlock()

	if s, exists := sessions[id]; exists {
		key = s.Key
		name = s.Name
		ok = true
		s.LastTime = time.Now()
	}
	//log.Info("get session id ", id, name, key, ok)
	return
}

//更新Session的过期时间
func FreshSession(id string) {
	lock.Lock()
	defer lock.Unlock()
	if s, ok := sessions[id]; ok {
		s.LastTime = time.Now()
	}
}

//删除一个Session
func RemoveSession(id string) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := sessions[id]; ok {
		delete(sessions, id)
	}
}
