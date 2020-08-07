package tdxF10Protocol_goVer

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strconv"
)

func (s *Socket) GetCompanyInfoContent(code, filename, start, length string) (error, string) {
	var maxRetry int
	if s.MaxRetry == 0 {
		maxRetry = len(s.Addrs)
	} else {
		maxRetry = s.MaxRetry
	}
	var err error
	for i := 0; i < maxRetry; i++ {
		err = s.NewConnectedSocket("")
		if err != nil {
			continue
		}
		err = s.setup()
		if err != nil {
			continue
		}

		if len(code) != 6 {
			return fmt.Errorf("corrupted code in category"), ""
		}
		if _, err := s.Client.Write(makeContentReq(code, filename, start, length)); err != nil {
			return err, ""
		}
		err, bodybuf := read(s.Client, s.Timeout)
		if err != nil {
			continue
		}
		err, content := parseContent(bodybuf)
		if err != nil {
			continue
		}

		return nil, content
	}
	return err, ""
}

func makeContentReq(code, filename, start, length string) []byte {
	//pkg.extend(struct.pack(u"<H6sH80sIII", market, code, 0, filename, start, length, 0))
	req, _ := hex.DecodeString("0c07109c000168006800d002") //header
	req = append(req, []byte{GetMarketByte(code), 0}...)   //market
	req = append(req, []byte(code)...)                     //code
	req = append(req, []byte{0, 0}...)                     //0
	req = append(req, []byte(filename)...)                 //filename
	req = append(req, make([]byte, 80-len([]byte(filename)))...)
	startI, _ := strconv.Atoi(start)
	lengthI, _ := strconv.Atoi(length)
	startb := make([]byte, 4)
	binary.LittleEndian.PutUint32(startb, uint32(startI)) //start
	req = append(req, startb...)
	lengthb := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthb, uint32(lengthI)) //length
	req = append(req, lengthb...)
	req = append(req, []byte{0, 0, 0, 0}...) //padding
	return req
}

func parseContent(bodybuf []byte) (error, string) {

	length := int(binary.LittleEndian.Uint16(bodybuf[10:12]))

	b, err := simplifiedchinese.GBK.NewDecoder().Bytes(bodybuf[12 : 12+length])
	return err, string(b)
}
