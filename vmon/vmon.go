package vmon

import (
	"../vcfg"
	"../vlog"
	"encoding/json"
	"fmt"
	"net/http"
)

type Counters struct {
	Remote   []string
	BytesIn  int64
	BytesOut int64
}

func (c *Counters) In(b int64) {
	c.BytesIn += b
}

func (c *Counters) Out(b int64) {
	c.BytesOut += b
}

var (
	counters map[string]Counters
)

func In(v *vcfg.Upstream, b int64) {
	c := counters[v.Local]
	c.In(b)
	counters[v.Local] = c
}

func Out(v *vcfg.Upstream, b int64) {
	c := counters[v.Local]
	c.Out(b)
	counters[v.Local] = c
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
	counters = make(map[string]Counters)
	for _, v := range *upstreams {
		counters[v.Local] = Counters{Remote: v.Remote}
		vlog.Info("%s %#v %v", v.Local, counters[v.Local], v.Remote)
	}
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":1972", nil)
	if err != nil {
		vlog.Err("binding error: ", err)
	}
}
