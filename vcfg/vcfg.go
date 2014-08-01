package vcfg

import (
	"encoding/json"
	"io/ioutil"
)

type Upstream struct {
	Local       string
	Remote      []string
	Connections []int64
	Timeout     int
	N           int
}

func ReadConfig(f *string) (*[]Upstream, error) {
	b, err := ioutil.ReadFile(*f)
	if err != nil {
		return nil, err
	}
	var upstreams []Upstream
	err = json.Unmarshal(b, &upstreams)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(upstreams); i++ {
		n := len(upstreams[i].Remote)
		upstreams[i].Connections = make([]int64, n, n)
	}
	return &upstreams, nil
}
