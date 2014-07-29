package vnet

import (
	"../vcfg"
	"../vlog"
	"../vmon"
	"../vutil"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	IN  = "<"
	OUT = ">"
)

func doForward(dir string, v *vcfg.Upstream, in, out net.Conn, p, src, dst net.Addr) {
	n, err := io.Copy(in, out)
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
		err = nil
	}
	if dir == IN {
		vmon.Inc(v, n, 0, err)
		out.Close()
		in.Close()
	} else if dir == OUT {
		vmon.Inc(v, 0, n, err)
		out.Close()
		in.Close()
	}
	vlog.Info("%v %s %v %v %v [%v]", p, dir, src, dst, n, vutil.T((err != nil), err, "OK"))
}

func goForward(local net.Conn, v vcfg.Upstream, i int) {
	dialer := net.Dialer{Timeout: time.Second * time.Duration(v.Timeout), KeepAlive: time.Second * time.Duration(v.KeepAlive)}
	remote, err := dialer.Dial("tcp", v.Remote[i])
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	vmon.AddStats()
	go doForward(OUT, &v, remote, local, local.LocalAddr(), remote.LocalAddr(), remote.RemoteAddr())
	go doForward(IN, &v, local, remote, local.LocalAddr(), remote.RemoteAddr(), remote.LocalAddr())
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
