package main

import (
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/tdxF10Protocol-goVer"
	"time"
)

func main() {
	// this mod subjects to codes that start with 000, 001, 002, 300 in sz and, 600, 601, 603 in sh
	// retry with each ip in the ipPool recursively; MaxRetry=len(Addrs) if 0/undefined
	api := tdxF10Protocol_goVer.Socket{
		MaxRetry: 0,
	}
	// Addrs affects GetCompanyInfoCategory & GetCompanyInfoContent; nil for default
	// timeout=5s by default(<=0)
	api.Init(nil, 0)
	defer api.Close()

	// get code-name mappings
	// it uses a separate ip pool; nil for default
	t := time.Now()
	err, codeNameMap := api.GetCodeNameMap(nil)
	if err != nil {
		println(err.Error())
	}

	fmt.Printf("codename len %d\n", len(codeNameMap))
	fmt.Printf("GetCodeNameMap took %s\n", time.Since(t))
	_ = PrettyPrint(codeNameMap)

	// get fields to fetch the f10 content of given codes
	// this method applies the retry policy
	codes := []string{}
	for k := range codeNameMap {
		codes = append(codes, k)
	}
	t = time.Now()
	err, category := api.GetCompanyInfoCategory(codes)
	if err != nil {
		fmt.Println(err.Error())
	}

	_ = PrettyPrint(category)
	fmt.Printf("category len %d\n", len(category))
	fmt.Printf("GetCompanyInfoCategory took %s\n", time.Since(t))

	// get the specific content
	code := "600000"
	x := category[code]["经营分析"]
	filename := x[0]
	start := x[1]
	length := x[2]
	t = time.Now()
	err, s := api.GetCompanyInfoContent(code, filename, start, length)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(s)
	fmt.Printf("GetCompanyInfoContent took %s", time.Since(t))

}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
