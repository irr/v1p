package vmon

import (
	"../vcfg"
	"../vlog"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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
	wg       sync.WaitGroup
	counters map[string]*Counters
)

func Inc(v *vcfg.Upstream, in, out int64, err error) {
	go func() {
		wg.Wait()
		wg.Add(1)
		c := counters[v.Local]
		c.Inc(in, out, err)
		defer wg.Done()
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

func Start(upstreams *[]vcfg.Upstream) {
	counters = make(map[string]*Counters)
	for _, v := range *upstreams {
		counters[v.Local] = &Counters{Remote: v.Remote}
	}
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":1972", nil)
	if err != nil {
		vlog.Err("binding error: ", err)
	}
}
