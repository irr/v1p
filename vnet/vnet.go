package vnet

import (
	"../vcfg"
	"../vlog"
	"../vutil"
	"io"
	"net"
	"time"
)

func doForward(in, out net.Conn) {
	n, err := io.Copy(in, out)
	vlog.Info("%v/%v -> %v/%v = %v bytes [%v]", in.LocalAddr(), in.RemoteAddr(), out.LocalAddr(), out.RemoteAddr(),
		n, vutil.T((err != nil), err, "OK"))
}

func Forward(local net.Conn, v vcfg.Upstream) {
	dialer := net.Dialer{Timeout: time.Second * time.Duration(v.Timeout), KeepAlive: time.Second * time.Duration(v.KeepAlive)}
	remote, err := dialer.Dial("tcp", *v.Remote)
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	go doForward(local, remote)
	go doForward(remote, local)
}
