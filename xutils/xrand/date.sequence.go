package xrand

import (
	"fmt"
	"sync"
	"time"
)

var pSingleton *_DateSequence
var once sync.Once

// 日期序列
type _DateSequence struct {
	Date     time.Time
	Sequence int
	lock     sync.Mutex
}

// 获取单例
func DateSequence_Singleton() *_DateSequence {
	return pSingleton
}

// 获取新序列
func (p *_DateSequence) ObtainSequnce() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	now := time.Now().Round(time.Second)
	if now.YearDay() > p.Date.YearDay() {
		p.Date = now
		p.Sequence = 0
	}

	// p.Sequence += random.Intn(10)
	p.Sequence++
	sequence := fmt.Sprintf("%0.3d", p.Sequence)
	return p.Date.Format("20060102-") + sequence
}

func init() {
	once.Do(func() { pSingleton = &_DateSequence{Date: time.Now().Round(time.Second)} })
}
