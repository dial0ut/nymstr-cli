package socks5

import (
    "net"
    "golang.org/x/net/proxy"
)

func DialSocks5(network, addr string) (net.Conn, error) {
    dialer, err := proxy.SOCKS5("tcp", "localhost:1080", nil, proxy.Direct)
    if err != nil {
        return nil, err
    }

    return dialer.Dial(network, addr)
}

