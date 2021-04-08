package tdxF10Protocol_goVer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var waditu = "http://api.waditu.com"

func GetCodeNameFromWaditu(token string) (error, map[string]string) {
	reqM := map[string]string{
		"api_name": "stock_basic",
		"token":    token,
		"fields":   "symbol,name",
	}
	type Data struct {
		Items [][]string `json:"items"`
	}
	type res struct {
		Data Data `json:"data"`
	}
	req, err := json.Marshal(reqM)
	if err != nil {
		return err, nil
	}
	rsp, err := http.Post(waditu, "text/plain", bytes.NewBuffer(req))
	if err != nil {
		return err, nil
	}
	defer rsp.Body.Close()
	var r res
	rb, _ := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(rb, &r)
	if err != nil {
		return err, nil
	}
	codename := map[string]string{}
	for _, pair := range r.Data.Items {
		if len(pair) != 2 {
			return fmt.Errorf("waditu rsp corrupted? %+v", pair), nil
		}
		if _, ok := preset[pair[0][:3]]; !ok {
			continue
		}
		codename[pair[0]] = pair[1]
	}
	return nil, codename
}
