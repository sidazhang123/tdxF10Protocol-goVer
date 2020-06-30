package main

import (
"bytes"
"compress/zlib"
"encoding/binary"
"encoding/hex"
"fmt"
"golang.org/x/text/encoding/simplifiedchinese"
"golang.org/x/text/transform"
"io"
"io/ioutil"
"net"
"strings"
)



func GetCompanyInfoCategory(addr string,code string)(error,[]map[string]interface{}){
	if len(code)!=6{return fmt.Errorf("corrupted code in category"),nil}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err!=nil{return err,nil}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err!=nil{return err,nil}
	defer conn.Close()
	for _,s:=range []string{"0c0218930001030003000d0001","0c0218940001030003000d0002","0c031899000120002000db0fd5d0c9ccd6a4a8af0000008fc22540130000d500c9ccbdf0d7ea00000002"}{
		if err:=setup(conn,s);err!=nil{return err,nil}
	}
	shSz := map[string]byte{"600": byte(1), "601": byte(1), "603": byte(1), "000": byte(0), "001": byte(0), "002": byte(0), "300": byte(0)}
	for k, market := range shSz {
		if strings.HasPrefix(code, k) {
			if _, err := conn.Write(makeCategoryReq(code,market));err!=nil{return err,nil}
			break
		}
	}
	err,bodybuf := read(conn)
	if err!=nil{return err,nil}
	return parseResponse(bodybuf)
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
func parseResponse(bodybuf []byte) (error,[]map[string]interface{}) {

	num := int(binary.LittleEndian.Uint16(bodybuf[:2]))
	if len(bodybuf)<num*152+2{return fmt.Errorf("category body len=%d while %d is required",len(bodybuf),num*152+2),nil}
	pos := 2
	category := []map[string]interface{}{}
	for i := 0; i <num; i++ {
		m:=make(map[string]interface{})
		nameb:=bodybuf[pos:pos+64]
		pos+=64
		filenameb:=bodybuf[pos:pos+80]
		pos+=80
		startb:=bodybuf[pos:pos+4]
		pos+=4
		lengthb:=bodybuf[pos:pos+4]
		pos+=4
		m["name"]=getStr(nameb)
		m["filename"]=getStr(filenameb)
		m["start"]=int(binary.LittleEndian.Uint32(startb))
		m["length"]=int(binary.LittleEndian.Uint32(lengthb))
		category=append(category,m)
	}
	return nil,category

}


func makeCategoryReq(code string,market byte) []byte {
	req, _ := hex.DecodeString("0c0f109b00010e000e00cf02") //header
	req = append(req, []byte{market, 0}...)       //market
	req = append(req, []byte(code)...)   //code
	req = append(req, []byte{0, 0, 0, 0}...) //padding
	return req
}

func setup(conn *net.TCPConn,hexStr string) error{
	req, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	_,err=conn.Write(req)
	if err != nil {
		return err
	}
	err,_=read(conn)
	if err != nil {
		return err
	}
	return nil
}

func read(conn *net.TCPConn) (error,[]byte) {
	//read header
	h := make([]byte, 16)
	_,err:=io.ReadFull(conn, h)
	if err != nil {
		return err,nil
	}
	//get zipsize&unzipsize

	unzipsize := int(binary.LittleEndian.Uint16(h[14:16]))
	zipsize := int(binary.LittleEndian.Uint16(h[12:14]))

	//read body
	body := []byte{}
	for {
		b := make([]byte, zipsize)
		_, e := io.ReadFull(conn, b)
		if e == nil {
			body = append(body, b...)
		}
		if e == io.EOF {
			body = append(body, b...)
			break
		} else {
			break
		}

	}
	//(unzip)
	if zipsize != unzipsize {
		debody, err := decompress(body)
		if err != nil {
			return err,nil
		}
		return nil,debody
	}
	return nil,body

}
func decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	p, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	return p, nil
}
