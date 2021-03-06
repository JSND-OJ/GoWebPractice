package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var userMap map[int]*User //메인함수 뉴핸들러 시점에서 초기화
var lastID int            //마지막 id 등록

// User struct
type User struct { //josn을 읽을 수 있는 스트럭트. id정수형 추가
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 //vars가 자동으로id 파싱
	id, err := strconv.Atoi(vars["id"]) //vars는 string이므로 atoi로 int정수형으로 바꾸면 첫번째 인티저형id와 두번쨰 err 반환
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id]
	if !ok { //ok가 없다면 즉 해당 id가 없으면
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(user) //go value로 된 유저정보를 json 문자열로 변환해서 data와 err 리턴, err는 무시
	fmt.Fprint(w, string(data))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// Created User
	lastID++ //id가 만들어질때마다 user르 기억하고 등록 .하나씩 증가
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user //id증가값 유저맵에 담김

	w.Header().Add("Content-Type", "application/json")
	data, _ := json.Marshal(user) //유저정보 마샬링(go value -> json)해서 바이트 어레이로 바꾸고 데이터에 넣어줌
	fmt.Fprint(w, string(data))
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User) //유저 스트럭트의 유저 정보를 담을 맵 생성
	lastID = 0                    //사용하기 전에 초기화
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)
	return mux
}
