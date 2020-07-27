package tdxF10Protocol_goVer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strconv"
)

/*
@ params [code1,code2,...]
@ return
{
	code:
		{name: [filename,start,length],
		 name: ,
		},
	code:,
}
*/

func (s *Socket) GetCompanyInfoCategory(codes, ipPool []string) (error, map[string]map[string][]string) {

	if ipPool == nil {
		ipPool = []string{}
	}
	err := s.NewConnectedSocket(ipPool)
	if err != nil {
		return err, nil
	}
	err = s.Setup()
	if err != nil {
		return err, nil
	}

	ret := map[string]map[string][]string{}
	for _, code := range codes {

		if len(code) != 6 {
			return fmt.Errorf("corrupted code in category"), nil
		}
		if _, err := s.Client.Write(makeCategoryReq(code)); err != nil {
			return err, nil
		}
		err, bodybuf := read(s.Client)

		if err != nil {
			return err, nil
		}
		err, category := parseCategory(bodybuf)
		if err != nil {
			return err, nil
		}
		ret[code] = category
	}
	return nil, ret
}

func makeCategoryReq(code string) []byte {
	req, _ := hex.DecodeString("0c0f109b00010e000e00cf02") //header
	req = append(req, []byte{GetMarketByte(code), 0}...)   //market
	req = append(req, []byte(code)...)                     //code
	req = append(req, []byte{0, 0, 0, 0}...)               //padding
	return req
}

func parseCategory(bodybuf []byte) (error, map[string][]string) {

	num := int(binary.LittleEndian.Uint16(bodybuf[:2]))
	if num == 0 {
		return fmt.Errorf("empty body, try another addr"), nil
	}
	if len(bodybuf) < num*152+2 {
		return fmt.Errorf("category body len=%d while %d is required", len(bodybuf), num*152+2), nil
	}
	pos := 2
	category := map[string][]string{}

	for i := 0; i < num; i++ {

		nameb := bodybuf[pos : pos+64]
		pos += 64
		filenameb := bodybuf[pos : pos+80]
		pos += 80
		startb := bodybuf[pos : pos+4]
		pos += 4
		lengthb := bodybuf[pos : pos+4]
		pos += 4

		category[getStr(nameb)] = []string{getStr(filenameb), strconv.Itoa(int(binary.LittleEndian.Uint32(startb))),
			strconv.Itoa(int(binary.LittleEndian.Uint32(lengthb)))}
	}
	fmt.Printf("%+v", category)
	return nil, category

}

func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func getStr(b []byte) string {
	p := bytes.IndexByte(b, byte(0))
	if p != -1 {
		b = b[0:p]
	}
	ub, e := gbkToUtf8(b)
	if e != nil {
		return "unknown_str"
	}
	return string(ub)
}
