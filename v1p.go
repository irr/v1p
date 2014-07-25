package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net"
	"os"
	"time"
)

const (
	T = "[v1p] "
)

var (
	L func(string, ...interface{})
	E func(string, ...interface{})
)

func forward(local net.Conn, remoteAddr *string, timeout int) {
	dialer := net.Dialer{Timeout: time.Second * time.Duration(timeout)}
	remote, err := dialer.Dial("tcp", *remoteAddr)
	if remote == nil {
		E("remote dial failed: %v", err)
		return
	}
	go io.Copy(local, remote)
	go io.Copy(remote, local)
}

func vip(l, r *string, t int) {
	L("proxying %s to %s (t:%d)...", *l, *r, t)
	local, err := net.Listen("tcp", *l)
	if local == nil {
		E("cannot listen: %v", err)
		os.Exit(1)
	}
	for {
		conn, err := local.Accept()
		if conn == nil {
			E("accept failed: %v", err)
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

	log.SetPrefix(T)

	if *s {
		syslog, err := syslog.New(syslog.LOG_INFO, T)
		if err != nil {
			log.Fatal(err)
		}
		L = func(f string, a ...interface{}) { syslog.Info(fmt.Sprintf(f, a...)) }
		E = func(f string, a ...interface{}) { syslog.Err(fmt.Sprintf(f, a...)) }
	} else {
		L = func(f string, a ...interface{}) { log.Println(fmt.Sprintf(f, a...)) }
		E = L
	}

	if *l != "" && *r != "" {
		vip(l, r, *t)
	} else {
		flag.Usage()
		os.Exit(0)
	}
}
