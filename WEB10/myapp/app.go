package myapp

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func NewHandler() http.Handler { //그냥 핸들러 만들어주는 친궂
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	return mux
}
