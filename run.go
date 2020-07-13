package main

import (
	"encoding/json"
	"fmt"
	"os"
	"tdxF10Protocol/te"
)

func main() {
	code := "600001"
	api := te.Socket{}
	defer api.Close()
	api.NewConnectedSocket([]string{})
	api.Setup()
	e, m := api.GetCompanyInfoCategory([]string{code})
	if e != nil {
		println(e.Error())
	}
	//PrettyPrint(m)
	x := m[code]["经营分析"]

	e, s := api.GetCompanyInfoContent(code, x[0], x[1], x[2])
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
	fmt.Println(s)

}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
