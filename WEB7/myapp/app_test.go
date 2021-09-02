package myapp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) //app.go 뉴핸들러 호출. 핸들러 만들어줘야해 왜? 반환을 처리
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jungehon", "last_name":"OH", "email":"cuttleoh@naver.com"}`)) // ``go가 이해할수 있도록 어노테이션. 안하면 문자열이 초기화->공백이 되어버려어엇
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//위 post방식 json정보를 서버가 user정보를 받아서 user정보를 리턴하는 부분
	user := new(User) //유저 맵 이라는 보관소

	err = json.NewDecoder(resp.Body).Decode(user) //서버가 보낸 데이터를 읽어온다, encoder/decoder는 스트림기반 데이터를 다루고, encoder는 go value를 json으로 반환
	assert.NoError(err)
	assert.NotEqual(0, user.ID) //유저 맵에 크리에이트된 값(user.ID)이 있다. user id 가 0이 아니다. 등록되어 있다

	id := user.ID                                               //             :=인수 입력. 함수호출시 함수로 값을 전달해주는 값
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id)) //id 값 순서확인. 타겟은 유저맵. 해당하는 유저맵 뒤져.무언가를 확인할 때는 겟방식. 무언가를 조회할려면 겟. Itoa 정수형을 문자열로 변환. id정보를 get 방식으로 user/id를 넣어서 오도록 만든다.get방식 쓰는 이유? 저장된 값들을 확인할려고
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	user2 := new(User)                             // :=인수 입력. 함수호출시 함수로 값을 전달해주는 값
	err = json.NewDecoder(resp.Body).Decode(user2) //decode는 {}를 읽어온다
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID) //user id는 create한 user.id이고, 새로운 정보를 받은 user2 id iser2 id
	assert.Equal(user.FirstName, user2.FirstName)

	log.Print(user)
	log.Print(user2.FirstName)
	log.Print(user2.LastName)
	log.Print(user2.Email)
	log.Print(user2.CreatedAt)

}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil) //delete는 기본 메소드가 아니다
	resp, err := http.DefaultClient.Do(req)                     //그래서 do로 위의 리퀘스트를 불러옴. Do메소드가 Delete메소드처럼 일함.
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	log.Print(string(data)) //
	assert.Contains(string(data), "No User ID:1")

	//만들고
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jungheon", "last_name":"Oh", "email":"cuttleoh@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	//다시 삭제
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User ID:1")
	log.Print(string(data))
}
