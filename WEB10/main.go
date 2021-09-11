package main

import (
	"GoBlock/goweb/WEB10/decoHandler"
	"GoBlock/goweb/WEB10/myapp"
	"log"
	"net/http"
	"time"
)

// logger = decoratorFunc 구현
func logger(w http.ResponseWriter, r *http.Request, h http.Handler) { //log데코레이터를 만든다, h http.handler 인자를 추가 받는다.w
	start := time.Now()                                                        // (4) req를 숳ㅇ하는데 걸리는 시간
	log.Println("[LOGGER1] Started")                                           // (1) req가 왔을 떄 handler 호출하기전에 먼저 started log를 찍는다
	h.ServeHTTP(w, r)                                                          // (2) http.ServeHTTP 호출
	log.Println("[LOGGER1] Completed time:", time.Since(start).Milliseconds()) // (3) 끝날 떄 logger 찍는다 (5) 특정 time으로부터 duration을 millisecnd 단위로
}

/*func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) { //log 데코레이터를 받는다
	start := time.Now()
	log.Println("[LOGGER2] Started") //handler 수행전 먼저 log를 찍는다
	h.ServeHTTP(w, r)                //요 핸들러는 log1
	log.Println("[LOGGER2] Completed time:", time.Since(start).Milliseconds())
}*/

func NewHandler() http.Handler {
	h := myapp.NewHandler()
	h = decoHandler.NewDecoHandler(h, logger) //로거1이 핸들러 가지고 있음. 새로운 ddecohandler패키지에서 myapp.newhandler() mux를 감싸준다(wrapping)
	//h = decoHandler.NewDecoHandler(h, logger2) //로거2가 로그1 가지고 있고
	return h
}

func main() {
	mux := NewHandler()

	http.ListenAndServe(":3004", mux)
}
