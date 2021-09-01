package main

import (
	"net/http"

	"GoBlock/goweb/WEB6/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}
