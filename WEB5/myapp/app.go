package myapp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World") //w에 hello world값을 주어서 프린트
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                   // r요청값을 vars가 받음.mux니까 라우터 역할.r은 request vars가 if 파싱을 자동으로 해준다.
	fmt.Fprint(w, "User Id:", vars["id"]) // vars["id"]형식으로 작성해야 파싱값이 user id로 들어가
}

func NewHandler() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler) //하위경로 미정의시 상위경로 자동호출
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)
	return mux
}
