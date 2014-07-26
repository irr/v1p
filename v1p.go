package main

import (
	"./vcfg"
	"./vlog"
	"./vnet"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
)

const (
	P = "[v1p] "
)

func vip(v vcfg.Upstream) {
	vlog.Info("proxying %s to %s (t:%d,k:%d)...", *v.Local, *v.Remote, v.Timeout, v.KeepAlive)
	local, err := net.Listen("tcp", *v.Local)
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
		go vnet.Forward(conn, v)
	}
}

func main() {
	s := flag.Bool("s", false, "syslog (enabled/disabled)")
	h := flag.Bool("h", false, "help")
	l := flag.String("l", "", "saddr:port (local)")
	r := flag.String("r", "", "raddr:port (remote)")
	c := flag.String("c", "", "config file")
	t := flag.Int("t", 0, "timeout (seconds)")
	k := flag.Int("k", 0, "keepalive (seconds)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"v1p version 0.1 (ivan.ribeiro@gmail.com)\nv1p [-s][-h][-t][-k] -l <addr:port> -r <addr:port>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	vlog.SetLogger(P, *s)

	barrier := make(chan int)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		<-sig
		barrier <- 1
	}()

	if *l != "" && *r != "" {
		upstream := vcfg.Upstream{Local: l, Remote: r, Timeout: *t, KeepAlive: *k}
		go vip(upstream)
		<-barrier
		os.Exit(0)
	} else if *c != "" {
		upstreams, err := vcfg.ReadConfig(c)
		if err != nil {
			vlog.Err("config error: %v", err)
		}
		for _, v := range *upstreams {
			go vip(v)
		}
		<-barrier
		os.Exit(0)
	} else {
		flag.Usage()
		os.Exit(0)
	}
}
