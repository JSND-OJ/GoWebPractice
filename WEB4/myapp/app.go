package myapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"` //어노테이션, 설명을 붙이는 것
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GoBlock to Busan")
}

type fooHandler struct{} //인스턴스 생성

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //servelHTTP 인터페이스 구현
	user := new(User)                           //post 전송방식은 a new foo를 created 반환.
	err := json.NewDecoder(r.Body).Decode(user) //string 기반 데이터전송, json을 go value로 스프링기반
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	user.CreatedAt = time.Now()

	data, _ := json.Marshal(user) //인터페이스. json형식 결과값은 바이트 어레이
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	log.Print(string(data))
	fmt.Fprint(w, string(data))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func NewHttpHandler() http.Handler {
	mux := http.NewServeMux() //라우터 클래스 만든다 mux 등록
	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", &fooHandler{})
	return mux
}
