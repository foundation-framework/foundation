package log

import (
	"net"
	"net/url"

	"go.uber.org/zap"
)

const (
	udpScheme = "udp"
)

type udpSink struct {
	conn *net.UDPConn
}

func newUdpSink(url *url.URL) (zap.Sink, error) {
	addr, err := net.ResolveUDPAddr("udp", url.Host)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &udpSink{conn: conn}, nil
}

func (u *udpSink) Write(buf []byte) (int, error) {
	return u.Write(buf)
}

func (u *udpSink) Sync() error {
	// Only for interface realization
	return nil
}

func (u *udpSink) Close() error {
	return u.conn.Close()
}
