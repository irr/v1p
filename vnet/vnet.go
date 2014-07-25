package vnet

import (
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

func Forward(local net.Conn, remoteAddr *string, timeout int) {
	dialer := vutil.T((timeout > 0), net.Dialer{Timeout: time.Second * time.Duration(timeout)}, net.Dialer{}).(net.Dialer)
	remote, err := dialer.Dial("tcp", *remoteAddr)
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	go doForward(local, remote)
	go doForward(remote, local)
}
