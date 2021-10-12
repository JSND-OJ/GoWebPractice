package main

import (
	"log"
	"net/http"

	"GoBlock/goweb/WEB20/app"

	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler("./member.db")
	defer m.Close()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
