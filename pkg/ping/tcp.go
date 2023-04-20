package ping

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

type TcpPingResult struct {
	Time int
	Err  error
	IP   net.IP
}

func (this *TcpPingResult) Result() int {
	return this.Time
}

func (this *TcpPingResult) Error() error {
	return this.Err
}

func (this *TcpPingResult) String() string {
	if this.Err != nil {
		return fmt.Sprintf("%s", this.Err)
	} else {
		return fmt.Sprintf("%s: time=%d ms", this.IP.String(), this.Time)
	}
}

type TcpPing struct {
	host    string
	Port    uint16
	Timeout time.Duration

	ip net.IP
}

func (this *TcpPing) SetHost(host string) {
	this.host = host
	this.ip = net.ParseIP(host)
}

func (this *TcpPing) Host() string {
	return this.host
}

func (this *TcpPing) Ping() IPingResult {
	return this.PingContext(context.Background())
}

func (this *TcpPing) PingContext(ctx context.Context) IPingResult {
	ip := cloneIP(this.ip)
	if ip == nil {
		var err error
		ip, err = LookupFunc(this.host)
		if err != nil {
			return &TcpPingResult{0, err, nil}
		}
	}
	dialer := &net.Dialer{
		Timeout:   this.Timeout,
		KeepAlive: -1,
	}
	t0 := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ip.String(), strconv.FormatUint(uint64(this.Port), 10)))
	if err != nil {
		return &TcpPingResult{0, err, nil}
	}
	defer conn.Close()
	return &TcpPingResult{int(time.Since(t0).Milliseconds()), nil, ip}
}

func NewTcpPing(host string, port uint16, timeout time.Duration) *TcpPing {
	return &TcpPing{
		host:    host,
		Port:    port,
		Timeout: timeout,
		ip:      net.ParseIP(host),
	}
}

var (
	_ IPing       = (*TcpPing)(nil)
	_ IPingResult = (*TcpPingResult)(nil)
)
