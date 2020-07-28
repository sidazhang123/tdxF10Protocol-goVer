package tdxF10Protocol_goVer

import (
	"net"
	"time"
)

type Socket struct {
	Client      *net.TCPConn
	Addrs       []string
	nextAddrIdx int
	MaxRetry    int
	Timeout     int
}

func (s *Socket) Init(addrs []string, timeout int) {
	if addrs == nil || len(addrs) == 0 {
		s.Addrs = []string{"218.75.75.20:7709", "211.100.23.200:7779", "58.49.110.76:7709", "101.71.255.135:7709", "211.100.23.202:7709"}
	} else {
		s.Addrs = addrs
	}
	if timeout <= 0 {
		s.Timeout = 5
	} else {
		s.Timeout = timeout
	}
}
func (s *Socket) NewConnectedSocket(ip string) (err error) {

	if s.Client != nil {
		if err := s.Client.Close(); err != nil {
			return err
		}
	}

	var tcpAddr *net.TCPAddr
	var nextAddrIdx int
	var addrs []string
	if len(ip) > 0 {
		nextAddrIdx = 0
		addrs = []string{ip}
	} else {
		nextAddrIdx = s.nextAddrIdx
		addrs = s.Addrs
	}
	for i := nextAddrIdx; i < len(addrs); i++ {
		tcpAddr, err = net.ResolveTCPAddr("tcp4", addrs[i])
		if i == len(s.Addrs)-1 && err != nil {
			return
		}
		if err == nil {
			for j := 0; j < 3; j++ {
				api, err := net.DialTCP("tcp", nil, tcpAddr)
				if err == nil {
					s.Client = api
					if len(ip) > 0 {
						s.nextAddrIdx = (i + 1) % len(s.Addrs)
					}
					return err
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	return
}

func (s *Socket) Close() error {
	return s.Client.Close()
}
