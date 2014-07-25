package main

import (
	"./vlog"
	"./vnet"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	P = "[v1p] "
)

func vip(l, r *string, t int) {
	vlog.Info("proxying %s to %s (t:%d)...", *l, *r, t)
	local, err := net.Listen("tcp", *l)
	if local == nil {
		vlog.Err("cannot listen: %v", err)
		os.Exit(1)
	}
	for {
		conn, err := local.Accept()
		if conn == nil {
			vlog.Err("accept failed: %v", err)
			os.Exit(1)
		}
		go vnet.Forward(conn, r, t)
	}
}

func main() {
	s := flag.Bool("s", false, "syslog (enabled/disabled)")
	h := flag.Bool("h", false, "help")
	l := flag.String("l", "", "saddr:port (local)")
	r := flag.String("r", "", "raddr:port (remote)")
	t := flag.Int("t", 0, "timeout (seconds)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "v1p version 0.1 (ivan.ribeiro@gmail.com)\n")
		fmt.Fprintf(os.Stderr, "v1p [-s][-h][-t] -l <addr:port> -r <addr:port>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	vlog.SetLogger(P, *s)

	if *l != "" && *r != "" {
		vip(l, r, *t)
	} else {
		flag.Usage()
		os.Exit(0)
	}
}
