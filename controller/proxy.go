package controller

import (
	"fmt"
	"net"
	"time"

	httpserver "github.com/Asutorufa/yuhaiin/net/proxy/http/server"
	proxyI "github.com/Asutorufa/yuhaiin/net/proxy/interface"
	"github.com/Asutorufa/yuhaiin/net/proxy/redir/redirserver"
	socks5server "github.com/Asutorufa/yuhaiin/net/proxy/socks5/server"
)

type sType int

var (
	hTTP   sType = 1
	socks5 sType = 2
	redir  sType = 3
	arr          = []sType{hTTP, socks5, redir}
)

type LocalListen struct {
	Server map[sType]proxyI.Server
	hosts  *Hosts
}

type Hosts struct {
	Socks5  string
	HTTP    string
	Redir   string
	TCPConn func(string) (net.Conn, error)
}

type LocalListenOption func(hosts *Hosts)

func NewLocalListenCon(option LocalListenOption) (l *LocalListen, err error) {
	l = &LocalListen{}

	if option == nil {
		return l, nil
	}
	hosts := &Hosts{
		TCPConn: func(s string) (net.Conn, error) {
			return net.DialTimeout("tcp", s, 5*time.Second)
		},
	}
	option(hosts)
	l.hosts = hosts
	if l.Server == nil {
		l.Server = map[sType]proxyI.Server{}
	}

	for index := range arr {
		l.Server[arr[index]] = l.newS(hosts, arr[index])
	}
	l.setTCPConn(hosts.TCPConn)
	return
}

func (l *LocalListen) SetAHost(option LocalListenOption) (erra error) {
	if option == nil {
		return nil
	}
	h := &Hosts{
		Redir:  l.hosts.Redir,
		HTTP:   l.hosts.HTTP,
		Socks5: l.hosts.Socks5,
	}
	option(h)
	for index := range arr {
		if l.Server[arr[index]] == nil {
			continue
		}
		err := l.Server[arr[index]].UpdateListen(l.getHost(h, arr[index]))
		if err != nil {
			erra = fmt.Errorf("%v\n UpdateListen %d -> %v", erra, arr[index], err)
		}
	}
	l.setTCPConn(h.TCPConn)
	l.hosts = h
	return
}

func (l *LocalListen) setTCPConn(conn func(string) (net.Conn, error)) {
	if conn == nil {
		return
	}
	for index := range arr {
		if l.Server[arr[index]] == nil {
			continue
		}
		l.Server[arr[index]].SetTCPConn(conn)
	}
}

func (l *LocalListen) newS(host *Hosts, sType2 sType) proxyI.Server {
	if host == nil {
		return nil
	}
	switch sType2 {
	case hTTP:
		server, _ := httpserver.New(host.HTTP)
		return server
	case socks5:
		server, _ := socks5server.New(host.Socks5)
		return server
	case redir:
		server, _ := redirserver.New(host.Redir)
		return server
	}
	return nil
}

func (l *LocalListen) getHost(option *Hosts, sType2 sType) string {
	if option == nil {
		return ""
	}
	switch sType2 {
	case hTTP:
		return option.HTTP
	case socks5:
		return option.Socks5
	case redir:
		return option.Redir
	}
	return ""
}
