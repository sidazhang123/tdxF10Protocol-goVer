package te

import (
	"net"
)

type Socket struct {
	Client *net.TCPConn
	Addrs  []string
}

func (s *Socket) NewConnectedSocket(addrs []string) (err error) {
	if s.Client != nil {
		s.Client.Close()
	}
	if len(addrs) == 0 {
		s.Addrs = []string{"218.75.75.20:7709", "211.100.23.200:7779", "58.49.110.76:7709", "101.71.255.135:7709", "211.100.23.202:7709"}
	} else {
		s.Addrs = addrs
	}
	var tcpAddr *net.TCPAddr
	for i, a := range s.Addrs {
		tcpAddr, err = net.ResolveTCPAddr("tcp4", a)
		if i == len(s.Addrs)-1 && err != nil {
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
