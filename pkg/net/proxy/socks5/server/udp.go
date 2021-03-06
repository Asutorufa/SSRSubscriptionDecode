package socks5server

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/net/utils"
)

// https://github.com/haxii/socks5/blob/bb9bca477f9b3ca36fa3b43e3127e3128da1c15b/udp.go#L20

func udpHandle(b []byte, f proxy.Proxy) ([]byte, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("normalHandleUDP() -> b byte array is empty")
	}
	/*
	* progress
	* 1. listener get client data
	* 2. get local/proxy packetConn
	* 3. write client data to local/proxy packetConn
	* 4. read data from local/proxy packetConn
	* 5. write data that from remote to client
	 */
	host, port, addrSize, err := ResolveAddr(b[3:])
	if err != nil {
		return nil, fmt.Errorf("resolve socks5 address failed: %v", err)
	}

	if net.ParseIP(host) == nil {
		addr, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			return nil, fmt.Errorf("resolve IP Addr failed: %v", err)
		}
		host = addr.IP.String()
	}

	h := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := f.PacketConn(h)
	if err != nil {
		return nil, fmt.Errorf("get packetConn from f failed: %v", err)
	}
	defer conn.Close()

	target, err := net.ResolveUDPAddr("udp", h)
	if err != nil {
		return nil, fmt.Errorf("resolve udp addr failed: %v", err)
	}
	_ = conn.SetDeadline(time.Now().Add(time.Second * 10))

	// write data to target and read the response back
	fmt.Println("UDP write", conn.LocalAddr(), "->", target)
	// fmt.Println("write data:", data, "origin:", b)
	_, err = conn.WriteTo(b[3+addrSize:], target)
	if err != nil {
		return nil, fmt.Errorf("write data to remote packetConn failed: %v", err)
	}

	respBuff := *utils.BuffPool.Get().(*[]byte)
	defer utils.BuffPool.Put(&(respBuff))
	// copy(respBuff[0:3], []byte{0, 0, 0})
	copy(respBuff[:3+addrSize], b[:3+addrSize]) // copy addr []byte{0,0,0,addr...}

	n, addr, err := conn.ReadFrom(respBuff[3+addrSize:])
	if err != nil {
		return nil, fmt.Errorf("read data From remote packetConn failed: %v", err)
	}
	fmt.Println("UDP read from", addr.String())
	// fmt.Println("read data", respBuff[:n])

	return respBuff[:n+3+addrSize], nil
}
