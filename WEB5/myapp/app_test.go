package myapp

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//1단계 TestIndex NewHandler() 등록빌드 검증테스팅
func TestIndex(t *testing.T) { //test pkg 가져다 쓸려면 파일명 _test, 인자로 testing pkg의 T스트럭트
	assert := assert.New(t) //어썰트 오브젝트 생성

	ts := httptest.NewServer(NewHandler()) //목업서버
	defer ts.Close()

	resp, err := http.Get(ts.URL)                //URL이 문자열로 들어감.test서버의 url을 넣고 리턴값은 response와 error
	assert.NoError(err)                          //에러가 엄써야 함
	assert.Equal(http.StatusOK, resp.StatusCode) //response의 statusCode가 스테이터스오케이(200)와 같아야 함.

	//2단계 hello world 데이터 테스팅 검사
	data, _ := ioutil.ReadAll(resp.Body) //ioutil로 body값을 모두 읽어오고, 리턴값은 바이트어레이형식 data와 error, error는 무시.
	log.Print(string(data))
	assert.Equal("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	assert.Contains(string(data), "Get UserInfo")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	assert.Contains(string(data), "User Id:89")

	resp, err = http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	assert.Contains(string(data), "User Id:56")

	resp, err = http.Get(ts.URL + "/users/10")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	assert.Contains(string(data), "User Id:10")

	resp, err = http.Get(ts.URL + "/users/999")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	log.Print(string(data))
	assert.Contains(string(data), "User Id:999")
}
