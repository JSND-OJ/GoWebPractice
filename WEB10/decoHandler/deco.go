package decoHandler

import "net/http"

type DecoratorFunc func(http.ResponseWriter, *http.Request, http.Handler) //3개의 리퀘스트 받는 펑션타입

type DecoHandler struct {
	fn DecoratorFunc
	h  http.Handler //이놈 자체가 http.handler 구현하고 있음
}

func (self *DecoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //decohandler가 http 데코레이터 핸들러 만든다
	self.fn(w, r, self.h) //self를 호출하고 wr 호출하고 self.h를 구현
}

func NewDecoHandler(h http.Handler, fn DecoratorFunc) http.Handler {
	return &DecoHandler{
		fn: fn,
		h:  h,
	}
}
