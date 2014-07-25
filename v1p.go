package main

import (
	"./vlog"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	P = "[v1p] "
)

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

func forward(local net.Conn, remoteAddr *string, timeout int) {
	dialer := T((timeout > 0), net.Dialer{Timeout: time.Second * time.Duration(timeout)}, net.Dialer{}).(net.Dialer)
	remote, err := dialer.Dial("tcp", *remoteAddr)
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	go io.Copy(local, remote)
	go io.Copy(remote, local)
}

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
		go forward(conn, r, t)
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
