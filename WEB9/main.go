package main

import (
	"GoBlock/goweb/WEB9/cipher"
	"GoBlock/goweb/WEB9/lzw"
	"fmt"
)

type Component interface {
	Operator(string)
}

var sentData string //전역변수, 아래 send data
var recvData string

//기본기능 본래기능 정의 단계. 기본기능? 데이터 전송!
type SendComponent struct{} //실제 기본기능

func (self *SendComponent) Operator(data string) { //실제 기본기능 구성후 Operator 호출
	//Send data
	sentData = data
}

// 데코레이터1: 압축 기능 구현
type ZipComponent struct { // 압축컴포넌트 데코레이터로 다른 컴포넌트 보관
	com Component // 데코레이터는 com 컴포넌트를 가진다
}

//data 압축하고 Component호출형태로 구성
func (self *ZipComponent) Operator(data string) { //com맴버, ZipComponent느 data압축 lzw부터 진행
	zipData, err := lzw.Write([]byte(data)) //Write는 string을 []byte로 바꿔주고 압축결과로 []byte, err
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(zipData)) //err가 없으면 본인(Self)의 Operator에 압축화한 data 호출. ZipComponent데코레이터가 데코레이트하고 있는 실제 컴포넌트의 self.com.Operator를 압축한 data를 호출
}

//데코레이터2: 암호화 기능 구현
type EncryptComponent struct { //Encrypt는 com뿐만 아니라 반드시 key깞을 가지는 데코레이터를 만든다
	key string
	com Component //데코레이터는 com Component를 가진다
}

func (self *EncryptComponent) Operator(data string) {
	encryptData, err := cipher.Encrypt([]byte(data), self.key) //key값 멤버로 가진다
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(encryptData)) //암호화된 데이터 컴포넌트
}

//여기까지는 본래기능 + 데코레이터1(압축기능) + 데코레이터2(함호화기능) -> sentData에 모두 응축

//암호 풀기
type DecryptComponent struct { //복호화 단계 컴포넌트, key값을 가진다
	key string
	com Component
}

func (self *DecryptComponent) Operator(data string) {
	decryptData, err := cipher.Decrypt([]byte(data), self.key) //암호가 풀린 데이터
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decryptData))
}

//압축 풀기
type UnzipComponent struct { //압축풀기 데이터 데코레이터
	com Component
}

func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

//recvData로 최종 응축
type ReadComponet struct{}

func (self *ReadComponet) Operator(data string) {
	recvData = data
}

//최종단계

func main() {
	//send component
	sender := &EncryptComponent{ //sender는 zipComponent화 SendComponent 모두 가진 상태로
		key: "asdf", //
		com:/*&ZipComponent{ //압호화된 컴포넌트는 ZipComponent를 가지고
		com: */&SendComponent{}, //ZipComponent는 SendComponent를 가진다
	} /*,*/
	/*}*/
	sender.Operator("Hello GoBlock") //sender 보내는 데이터 ""
	fmt.Println(sentData)            //암호화+압축 데이터는 최족적으로 sentData에 들어있음

	////receive component
	receiver := /*&UnzipComponent{
		com: */&DecryptComponent{
			key: "asdf",
			com: &ReadComponet{},
		} /*,*/
	/*}*/
	receiver.Operator(sentData)
	fmt.Println(recvData) //복호화 하고 압축을 푼 데이터 출력
}
