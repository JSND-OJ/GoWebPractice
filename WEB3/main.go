package main

import (
	"GoBlock/goweb/WEB3/myapp"
	"net/http"
)

func main() {

	http.ListenAndServe(":3000", myapp.NewHttpHandler()) //mux대신 NewHttphandler()
}
