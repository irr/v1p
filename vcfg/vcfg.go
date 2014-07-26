package vcfg

import (
	"encoding/json"
	"io/ioutil"
)

type Upstream struct {
	Local   *string
	Remote  *string
	Timeout *int
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
	return &upstreams, nil
}
