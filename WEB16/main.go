package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) { // formvalue에 msg와 name이 들어오면 sendmessage에 넣어줍니다
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("add user: %s", username))
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("left user: %s", username))
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message //채널 인스턴스 변수 선언

func sendMessage(name, msg string) { // push
	// send message to every clients
	msgCh <- Message{name, msg} // message 구조체에 name이랑 msg 담고 채널로 보냅니다.
}

func processMsgCh(es eventsource.EventSource) { // pop
	for msg := range msgCh { // 채널에 데이터가 들어오면 msg로 넘겨줍니다
		data, _ := json.Marshal(msg)                                                 // 마샬링하고
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond())) //이벤트소스로 전체 유저에세 메세지 보냅니다
	}
}

func main() {
	msgCh = make(chan Message)      // 채널 초기화
	es := eventsource.New(nil, nil) //이벤트소스 불러옵니다
	defer es.Close()                // 자원 끌어쓰니까 나중에 받아줍니다

	go processMsgCh(es) //실시간으로 채팅이 이루어지게 스레드돌립니다

	mux := pat.New() // 핸들러 전환은 자바스크립트가?
	mux.Post("/messages", postMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
