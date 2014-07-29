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

var (
	mutex    sync.Mutex
	counters map[string]*Counters
	stats    vutil.CAPArray
)

func Inc(v *vcfg.Upstream, in, out int64, err error) {
	go func() {
		mutex.Lock()
		c := counters[v.Local]
		c.Inc(in, out, err)
		mutex.Unlock()
	}()
}

func Server(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(counters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "%s\n", b)
	}
}

func Scheduler(c <-chan time.Time) {
	for now := range c {
		mutex.Lock()
		fmt.Printf("%v %#v\n", now, stats)
		mutex.Unlock()
	}
}

func Start(upstreams *[]vcfg.Upstream) {
	counters = make(map[string]*Counters)
	for _, v := range *upstreams {
		counters[v.Local] = &Counters{Remote: v.Remote}
	}
	stats = vutil.CAPArray{N: 60}
	stats.Fill(0)
	scheduler := time.Tick(1 * time.Second)
	go Scheduler(scheduler)
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":1972", nil)
	if err != nil {
		vlog.Err("binding error: ", err)
	}
}
