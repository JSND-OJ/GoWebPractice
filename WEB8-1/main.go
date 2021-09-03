package main

import (
	"net/http"

	"GoBlock/goweb/WEB8-1/myapp"
)

func main() {
	http.ListenAndServe(":3003", myapp.NewHandler())
}
