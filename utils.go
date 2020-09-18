package tdxF10Protocol_goVer

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func (s *Socket) setup() (err error) {
	for _, setupStr := range []string{"0c0218930001030003000d0001", "0c0218940001030003000d0002", "0c031899000120002000db0fd5d0c9ccd6a4a8af0000008fc22540130000d500c9ccbdf0d7ea00000002"} {
		if err := _setup(s.Client, setupStr, s.Timeout); err != nil {
			return err
		}
	}
	return nil
}

func GetMarketByte(code string) byte {
	shSz := map[string]byte{"600": byte(1), "601": byte(1), "603": byte(1), "605": byte(1), "000": byte(0), "001": byte(0), "002": byte(0), "003": byte(0), "300": byte(0)}
	for k, market := range shSz {
		if strings.HasPrefix(code, k) {
			return market
		}
	}
	return 9
}

func _setup(conn *net.TCPConn, hexStr string, timeout int) error {
	req, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	_, err = conn.Write(req)
	if err != nil {
		return err
	}
	err, _ = read(conn, timeout)
	if err != nil {
		return err
	}
	return nil
}

func read(conn *net.TCPConn, timeout int) (error, []byte) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	//read header
	h := make([]byte, 16)
	n, err := io.ReadFull(conn, h)
	if err != nil || n != 16 {
		return err, nil
	}
	//get zipsize&unzipsize

	unzipsize := int(binary.LittleEndian.Uint16(h[14:16]))
	zipsize := int(binary.LittleEndian.Uint16(h[12:14]))

	//read body
	var body []byte
	for {
		b := make([]byte, zipsize)
		_, e := io.ReadFull(conn, b)
		if e == nil {
			body = append(body, b...)
		}
		if e == io.EOF {
			body = append(body, b...)
			break
		} else {
			break
		}
	}
	//(unzip)
	if zipsize != unzipsize {
		debody, err := decompress(body)
		if err != nil {
			return err, nil
		}
		return nil, debody
	}
	return nil, body

}
func decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	if err := z.Close(); err != nil {
		return nil, err
	}
	return p, nil
}
