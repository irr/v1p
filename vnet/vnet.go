package vnet

import (
	"../vlog"
	"../vutil"
	"io"
	"net"
	"time"
)

func Forward(local net.Conn, remoteAddr *string, timeout int) {
	dialer := vutil.T((timeout > 0), net.Dialer{Timeout: time.Second * time.Duration(timeout)}, net.Dialer{}).(net.Dialer)
	remote, err := dialer.Dial("tcp", *remoteAddr)
	if remote == nil {
		vlog.Err("remote dial failed: %v", err)
		return
	}
	go io.Copy(local, remote)
	go io.Copy(remote, local)
}
