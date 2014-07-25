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
	L = "[v1p] "
)

var (
	F    = fmt.Sprintf
	logi func(string) error
	loge func(string) error
)

func forward(local net.Conn, remoteAddr *string, timeout int) {
	dialer := net.Dialer{Timeout: time.Second * time.Duration(timeout)}
	remote, err := dialer.Dial("tcp", *remoteAddr)
	if remote == nil {
		loge(F("remote dial failed: %v", err))
		return
	}
	go io.Copy(local, remote)
	go io.Copy(remote, local)
}

func vip(l, r *string, t int) {
	logi(F("proxying %s to %s (t:%d)...", *l, *r, t))
	local, err := net.Listen("tcp", *l)
	if local == nil {
		loge(F("cannot listen: %v", err))
		os.Exit(1)
	}
	for {
		conn, err := local.Accept()
		if conn == nil {
			loge(F("accept failed: %v", err))
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
		fmt.Fprintf(os.Stderr, "usage: v1p [-s][-h][-t] -l <addr:port> -r <addr:port>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	log.SetPrefix(L)

	if *s {
		syslog, err := syslog.New(syslog.LOG_INFO, L)
		if err != nil {
			log.Fatal("v1p: ", err)
		}
		logi = syslog.Info
		loge = syslog.Err
	} else {
		logi = func(s string) (err error) { log.Println(s); return }
		loge = func(s string) (err error) { log.Println(s); return }
	}

	if *l != "" && *r != "" {
		vip(l, r, *t)
	} else {
		flag.Usage()
		os.Exit(0)
	}
}
