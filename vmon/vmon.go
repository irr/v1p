package vmon

import (
	"../vcfg"
	"../vlog"
	"../vutil"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Counters struct {
	Remote   []string
	BytesIn  int64
	BytesOut int64
	Success  uint64
	Errors   uint64
}

type Stats struct {
	Date     time.Time
	Counters *map[string]*Counters
}

const (
	GAP_SECS   = 60
	GAP_LENGTH = 15
)

var (
	mutex     sync.Mutex
	upstreams *[]vcfg.Upstream
	stats     *vutil.CAPArray
)

func (c *Counters) Inc(in, out int64, err error) {
	if in > 0 {
		c.BytesIn += in
	}
	if out > 0 {
		c.BytesOut += out
	}
	if err == nil {
		c.Success += 1
	} else {
		c.Errors += 1
	}
}

func Inc(v *vcfg.Upstream, in, out int64, err error) {
	go func() {
		mutex.Lock()
		e, err := stats.Geth(1)
		if err == nil {
			m := e.(*map[string]*Counters)
			c := (*m)[v.Local]
			c.Inc(in, out, err)
		}
		mutex.Unlock()
	}()
}

func Server(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	ts := make([]Stats, stats.N, stats.N)
	now := time.Now()
	for i := 0; i < stats.N; i++ {
		v, _ := stats.Geth(i + 1)
		if v != nil {
			c := v.(*map[string]*Counters)
			ts[i] = Stats{Date: now, Counters: c}
		}
		now = now.Add(-1 * GAP_SECS * time.Second)
	}
	b, err := json.Marshal(ts)
	mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "%s\n", b)
	}
}

func BuildCounter() *map[string]*Counters {
	c := make(map[string]*Counters)
	for _, v := range *upstreams {
		c[v.Local] = new(Counters)
		c[v.Local].Remote = v.Remote
	}
	return &c
}

func Scheduler(c <-chan time.Time) {
	for {
		<-c
		mutex.Lock()
		stats.Push(BuildCounter())
		mutex.Unlock()
	}
}

func Start(ups *[]vcfg.Upstream) {
	upstreams = ups
	mutex.Lock()
	stats = &vutil.CAPArray{N: GAP_LENGTH}
	for i := 0; i < stats.N; i++ {
		stats.Push(BuildCounter())
	}
	mutex.Unlock()
	scheduler := time.Tick(GAP_SECS * time.Second)
	go Scheduler(scheduler)
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":1972", nil)
	if err != nil {
		vlog.Err("binding error: ", err)
	}
}
