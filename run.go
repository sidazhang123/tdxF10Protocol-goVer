package main

import (
	"encoding/json"
	"fmt"
	"tdxF10Protocol/te"
)

func main() {
	api := te.Socket{}
	api.NewConnectedSocket([]string{})
	api.Setup()
	_, m := api.GetCompanyInfoCategory([]string{"000001", "000002"})
	api.Close()
	PrettyPrint(m)
}
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
