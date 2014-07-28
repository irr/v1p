package vmon

import (
	"../vcfg"
	"../vlog"
	"fmt"
	"net/http"
)

func Server(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", "OK")
}

func Start(upstreams *[]vcfg.Upstream) {
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":1972", nil)
	if err != nil {
		vlog.Err("binding error: ", err)
	}
}
