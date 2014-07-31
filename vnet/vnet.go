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

func goForward(local net.Conn, v *vcfg.Upstream) {
	ok := false
	dialer := net.Dialer{Timeout: time.Second * time.Duration(v.Timeout)}
	for i := 0; i < len(v.Remote); i++ {
		v.N = vutil.T((v.N >= len(v.Remote)), 0, v.N).(int)
		remote, err := dialer.Dial("tcp", v.Remote[v.N])
		if err == nil {
			v.N = v.N + 1
			go doForward(OUT, v, remote, local, local.LocalAddr(), remote.LocalAddr(), remote.RemoteAddr())
			go doForward(IN, v, local, remote, local.LocalAddr(), remote.RemoteAddr(), remote.LocalAddr())
			ok = true
			break
		}
		v.N = v.N + 1
	}
	if !ok {
		vlog.Err("%v %s %v %v [%v]", local.LocalAddr(), OUT, v.Remote, 0, "ERR")
		local.Close()
	}
}

func Vip(v *vcfg.Upstream) {
	vlog.Info("proxying %s to %v (t:%d)...", v.Local, v.Remote, v.Timeout)
	local, err := net.Listen("tcp", v.Local)
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
		go goForward(conn, v)
	}
}
