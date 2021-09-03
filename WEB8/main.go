package main

import (
	"net/http"

	"GoBlock/goweb/WEB8/myapp"
)

func main() {
	http.ListenAndServe(":3003", myapp.NewHandler())
}
