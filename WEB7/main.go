package main

import (
	"net/http"

	"GoBlock/goweb/WEB7/myapp"
)

func main() {
	http.ListenAndServe(":3002", myapp.NewHandler())
}
