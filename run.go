package main

import (
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/tdxF10Protocol-goVer/te"
)

func main() {
	// this mod subjects to codes that start with 000, 001, 002, 300, 600, 601, 603
	api := te.Socket{}
	defer api.Close()
	// get a code-name mapping
	err, codeNameMap := api.GetCodeNameMap(nil)
	if err != nil {
		println(err.Error())
	}
	fmt.Println(len(codeNameMap))
	//PrettyPrint(codeNameMap)
	// get fields to fetch the f10 content of given codes
	//code := "600001"
	//err, codeF10InfoMap := api.GetCompanyInfoCategory([]string{code},nil)
	//if err != nil {
	//	println(err.Error())
	//}
	//PrettyPrint(codeF10InfoMap)
	//// get the specific content
	//x := codeF10InfoMap[code]["经营分析"]
	//err, s := api.GetCompanyInfoContent(code, x[0], x[1], x[2],nil)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}
	//fmt.Println(s)

}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
