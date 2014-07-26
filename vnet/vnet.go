package vnet

import (
	"../vcfg"
	"../vlog"
	"../vutil"
	"io"
	"net"
	"os"
	"time"
)

func doForward(in, out net.Conn) {
	n, err := io.Copy(in, out)
	vlog.Info("%v/%v -> %v/%v = %v bytes [%v]", in.LocalAddr(), in.RemoteAddr(), out.LocalAddr(), out.RemoteAddr(),
		n, vutil.T((err != nil), err, "OK"))
}

func goForward(local net.Conn, v vcfg.Upstream, i int) {
	dialer := net.Dialer{Timeout: time.Second * time.Duration(v.Timeout), KeepAlive: time.Second * time.Duration(v.KeepAlive)}
	remote, err := dialer.Dial("tcp", v.Remote[i])
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	go doForward(local, remote)
	go doForward(remote, local)
}

func Vip(v vcfg.Upstream) {
	vlog.Info("proxying %s to %v (t:%d,k:%d)...", v.Local, v.Remote, v.Timeout, v.KeepAlive)
	local, err := net.Listen("tcp", v.Local)
	if local == nil {
		vlog.Err("cannot listen: %v", err)
		os.Exit(1)
	}
	n := 0
	for {
		conn, err := local.Accept()
		if conn == nil {
			vlog.Err("accept failed: %v", err)
			os.Exit(1)
		}
		n = vutil.T((n >= len(v.Remote)), 0, n).(int)
		go goForward(conn, v, n)
		n = n + 1
	}
}
