package te

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"sort"
)

// derived from get_security_list() in pytdx

func (s *Socket) GetCodeNameMap(ipPool []string) (error, map[string]string) {
	/*
		Code-Name lists are not always the same in different addrs. So pick a pool.
		Because it is patently fast, get the lists from all the addrs and return the longest.
	*/
	if ipPool == nil {
		ipPool = []string{"211.100.23.200:7779", "58.49.110.76:7709", "101.71.255.135:7709", "211.100.23.202:7709", "218.75.75.20:7709"}
	}
	ret := []map[string]string{}
	for _, ip := range ipPool {
		_ = newSocket(s, ip)
		codeNameMap := map[string]string{}
		fmt.Println(ip)
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
					fmt.Println(pos)
					break
				} else {
					extendMap(codeNameMap, codeName)
					pos += len(codeName)
				}

			}
			//if market==0{fmt.Println("hey")}
		}
		ret = append(ret, codeNameMap)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return len(ret[i]) > len(ret[j])
	})
	return nil, ret[0]
}

func extendMap(o, n map[string]string) {
	for k, v := range n {
		o[k] = v
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
	req = append(req, []byte{byte(market), 0}...)
	startb := make([]byte, 2)
	binary.LittleEndian.PutUint16(startb, uint16(start)) //start
	req = append(req, startb...)
	return req
}

func parseCodename(bodybuf []byte) (error, map[string]string) {

	num := int(binary.LittleEndian.Uint16(bodybuf[:2]))
	println(num)
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
		if string([]rune(code)[0]) == "6" {
			println(code)
		}
		if GetMarketByte(code) != 9 {
			codename[code] = string(name)
		}

	}
	return nil, codename

}
