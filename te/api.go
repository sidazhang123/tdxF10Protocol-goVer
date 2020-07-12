package te

import (
	"net"
)

type Socket struct {
	Client *net.TCPConn
}

func (s *Socket) NewConnectedSocket(addrs []string) (err error) {
	if len(addrs) == 0 {
		addrs = append(addrs, "211.100.23.200:7779")
	}
	var tcpAddr *net.TCPAddr
	for i, a := range addrs {
		tcpAddr, err = net.ResolveTCPAddr("tcp4", a)
		if i == len(addrs)-1 && err != nil {
			return
		}
		if err == nil {
			for i := 0; i < 3; i++ {
				api, err := net.DialTCP("tcp", nil, tcpAddr)
				if err == nil {
					s.Client = api
					return err
				}
			}
		}
	}
	return
}

func (s *Socket) Close() {
	s.Client.Close()
}
