package main

import "fmt"

func main() {
	//211.100.23.200:7779

	err,m:=GetCompanyInfoCategory("211.100.23.200:7779","000001")
	if err!=nil{fmt.Println(err.Error())}
	fmt.Printf("%+v",m)
}