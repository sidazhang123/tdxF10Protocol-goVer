package tdxF10Protocol_goVer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"sort"
	"strings"
)

// derived from get_security_list() in pytdx

func (s *Socket) GetCodeNameMap(ipPool []string) (error, map[string]string) {
	/*
		Code-Name lists are not always the same in different addrs. So pick a pool.
		Because it is patently fast, get the lists from all the addrs and return the longest.
	*/
	if ipPool == nil {
		ipPool = []string{"211.100.23.200:7779", "58.49.110.76:7709", "211.100.23.202:7709"}
	}
	var ret []map[string]string
	for _, ip := range ipPool {
		_ = newSocket(s, ip)
		codeNameMap := map[string]string{}
		for market := 0; market < 2; market++ {
			pos := 0
			for {
				if _, err := s.Client.Write(makeCodenameReq(market, pos)); err != nil {
					return err, nil
				}
				err, bodybuf := read(s.Client)
				if err != nil {
					return err, nil
				}
				err, codeName := parseCodename(bodybuf)
				if err != nil {
					return err, nil
				}
				if len(codeName) == 0 {
					break
				} else {
					pos += len(codeName)
				}
				extendMap(codeNameMap, codeName, market)
			}
		}
		ret = append(ret, codeNameMap)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return len(ret[i]) > len(ret[j])
	})
	return nil, ret[0]
}
func isInterested(market int, code string) bool {
	shSz := map[string]int{"600": 1, "601": 1, "603": 1, "000": 0, "001": 0, "002": 0, "300": 0}
	for preCode, preMarket := range shSz {
		if strings.HasPrefix(code, preCode) && preMarket == market {
			return true
		}
	}
	return false
}
func extendMap(o, n map[string]string, market int) {
	for k, v := range n {
		if isInterested(market, k) {
			o[k] = strings.Replace(v, " ", "", -1)
		}
	}
}
func newSocket(s *Socket, ip string) error {
	err := s.NewConnectedSocket([]string{ip})
	if err != nil {
		return err
	}
	err = s.Setup()
	if err != nil {
		return err
	}
	return nil
}

func makeCodenameReq(market, start int) []byte {

	//pkg.extend(struct.pack(u"<H6sH80sIII", market, code, 0, filename, start, length, 0))
	req, _ := hex.DecodeString("0c0118640101060006005004") //header
	marketb := make([]byte, 2)
	binary.LittleEndian.PutUint16(marketb, uint16(market)) //market
	req = append(req, []byte{byte(market), 0}...)
	startb := make([]byte, 2)
	binary.LittleEndian.PutUint16(startb, uint16(start)) //start
	req = append(req, startb...)
	return req
}

func parseCodename(bodybuf []byte) (error, map[string]string) {

	num := int(binary.LittleEndian.Uint16(bodybuf[:2]))

	pos := 2
	codename := map[string]string{}
	for i := 0; i < num; i++ {

		code := string(bodybuf[pos : pos+6])
		nameb := bytes.Trim(bodybuf[pos+8:pos+16], "\x00")
		name, err := simplifiedchinese.GBK.NewDecoder().Bytes(nameb)
		if err != nil {
			return err, nil
		}
		pos += 29
		codename[code] = string(name)
	}
	return nil, codename

}
