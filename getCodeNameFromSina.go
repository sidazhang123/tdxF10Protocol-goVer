package tdxF10Protocol_goVer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// this method is not part of the pytdx but it provides a more stable way to get the mapping
// typically, it uses yesterday's real data before 13:00 CST and updates to the latest after
// while the pytdx approach seems do not update sooner than the market closes.

/*
get
pagination
re code,name fields
gbk
return map[string]string  :: with no
*/
var sina = "http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData?num=200&sort=code&asc=0&node=%s&symbol=&_s_r_a=page&page=%d"
var preset = map[string]bool{"600": false, "601": false, "603": false, "605": false, "000": false, "001": false, "002": false, "003": false, "300": false}

func u2s(from string) (to string, err error) {
	a := regexp.MustCompile(`^([^\\]*)\\u`).FindStringSubmatch(from)
	prefix := ""
	if len(a[1]) > 0 {
		prefix = a[1]
		from = strings.TrimLeft(from, prefix)
	}
	suffix := ""
	b := regexp.MustCompile(`\\u[\da-f]{4}(.?)$`).FindStringSubmatch(from)
	if len(b[1]) > 0 {
		suffix = b[1]
		from = strings.TrimRight(from, suffix)
	}
	from = strings.ReplaceAll(from, " ", "")
	bs, err := hex.DecodeString(strings.Replace(from, `\u`, ``, -1))
	if err != nil {
		println(fmt.Sprintf("%+v\n", from))
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		_ = binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}

	to = prefix + to + suffix
	to = strings.ReplaceAll(to, " ", "")
	return
}
func GetCodeNameFromSina() (error, map[string]string) {
	var codename = make(map[string]string)
	var _url string
	for page := 0; page < 60; page++ {

		if page == 0 {
			_url = fmt.Sprintf(sina, "shfxjs", 1)
		} else {
			_url = fmt.Sprintf(sina, "hs_a", page)
		}
		resp, err := http.Get(_url)
		if err != nil || resp == nil {
			return err, codename
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err, codename
		}

		matches := regexp.MustCompile(`"code":"([\d]{6})","name":"([^"]+)"`).FindAllStringSubmatch(string(b), -1)
		if matches == nil {
			return nil, codename
		}
		for _, match := range matches {

			code := match[1]

			if _, ok := preset[code[:3]]; !ok {
				continue
			}
			name, err := u2s(match[2])
			if err != nil {
				return err, codename
			}
			codename[code] = name
		}
	}

	return nil, codename
}
