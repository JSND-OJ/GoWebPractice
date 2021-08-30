package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//--index Handler res.body로 render된 data를 비교검증하는 부분 //
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()               //실제 response를 사용하지 않고 테스팅 응신방법
	req := httptest.NewRequest("GET", "/", nil) //전송팡식 지정(get) 목적지(타겟)을 파싱
	//mux로 "/"경로를 타겟 라우팅(분배) 응신 렌더링 시키는 부분
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req) //인터페이스 왔다갔다 외부공개 기능

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body) //hello world 문자열 읽어옴여
	//app.go페이지에서 index Handler에서 render data를 읽어오는 부분---//
	log.Print(string(data))
	assert.Equal("GoBlock to Busan", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	log.Print(string(data))
	assert.Equal("Hello World!", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=O.J", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	log.Print(string(data))
	assert.Equal("Hello O.J!", string(data))
}

func TestFooPathHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	data, _ := ioutil.ReadAll(res.Body)
	log.Print(string(data))
}

func TestFooPathHandler_WithJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{"first_name":"J", "last_name":"O", "email":"cuttleoh@naver.com"}`))

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)                             //user 구조체 생성
	err := json.NewDecoder(res.Body).Decode(user) // json을 go value로 변환 스트림 방식
	assert.Nil(err)                               //객체가 비었는지 확인
	assert.Equal("J", user.FirstName)
	assert.Equal("O", user.LastName)

}
