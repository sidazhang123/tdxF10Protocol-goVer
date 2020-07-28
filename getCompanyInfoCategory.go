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

func (s *Socket) GetCompanyInfoCategory(codeSlice []string) (error, map[string]map[string][]string) {
	// handy for removal
	codes := map[string]bool{}
	for _, c := range codeSlice {
		codes[c] = true
	}

	ret := map[string]map[string][]string{}

	var bodybuf []byte
	var err error
	var category map[string][]string
	// try with each of the ipPool for one time; return
	var maxRetry int
	if s.MaxRetry == 0 {
		maxRetry = len(s.Addrs)
	} else {
		maxRetry = s.MaxRetry
	}
	for i := 0; i < maxRetry; i++ {
		var success []string
		err = s.NewConnectedSocket("")
		if err != nil {
			fmt.Println("new sock err")
			continue
		}
		err = s.setup()
		if err != nil {
			fmt.Println("setup err")
			continue
		}
		for code := range codes {
			if len(code) != 6 {
				return fmt.Errorf("corrupted code in category"), nil
			}

			_, err = s.Client.Write(makeCategoryReq(code))
			if err != nil {
				continue
			}
			err, bodybuf = read(s.Client, s.Timeout)
			if err != nil {
				continue
			}
			err, category = parseCategory(bodybuf)
			if err != nil {
				continue
			}
			ret[code] = category
			success = append(success, code)
		}

		for _, f := range success {
			delete(codes, f)
		}
		if len(codes) == 0 {
			break
		}

	}
	if len(codes) != 0 {
		ec := ""
		for k := range codes {
			ec += k + ","
		}
		err = fmt.Errorf("getCompanyInfoCategory err in %s", ec)
	}
	return err, ret
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
