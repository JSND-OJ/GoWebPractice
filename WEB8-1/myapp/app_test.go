package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
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
	assert.Equal(string(data), "No Users")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:89")

	resp, err = http.Get(ts.URL + "/users/9999")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:9999")
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jungheon", "last_name":"Oh", "email":"cuttleoh@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
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

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jungheon", "last_name":"Oh", "email":"cuttleoh@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User ID:1")
	log.Print(string(data))
}

func TestUpdateUser(t *testing.T) { //put전송방식은 update로 반환받는다
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"id":1, "first_name":"updated", "last_name":"updated", "email":"updated@naver.com"}`))
	resp, err := http.DefaultClient.Do(req) //do(req): do(method)="put" 메소드를 정의
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	//위에서 확인한 결과 id가 없으므로 post방식을 써서 만들어준다
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jungheon", "last_name":"Oh", "email":"cuttleoh@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user) //서버가 보낸 정보를 읽어온다. 매개변수 입력
	assert.NoError(err)
	assert.NotEqual(0, user.ID)
	log.Print("FistName Before update:", user.FirstName)

	//
	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"Bin"}`, user.ID)
	//%d는 동적 id값을 넣을 수 있도록 정리하고 update요청(패킷)시에 first_name만 셋팅해서 보냈다
	//first_name만 보내면 user구조체 json에서 Last_Name과 Email은 string형의 기본값, 즉 "" 공백문자열만 들어간다

	req, _ = http.NewRequest("PUT", ts.URL+"/users", //put메소드는 기본 메소드가 아니므로 do(req) request에 메소드로 정의
		strings.NewReader(updateStr)) //io로 updateStr 받는다
	resp, err = http.DefaultClient.Do(req) //do(req): do(method)="put" 메소드를 정의
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)                             //서버가 upadte된 user정보 처리
	err = json.NewDecoder(resp.Body).Decode(updateUser) //서버가 보낸 데이터를 읽어온다 = 매개변수 입력
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("Bin", updateUser.FirstName)
	assert.Equal(user.LastName, updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)

	log.Print(user)
	log.Print("User ID:", updateUser.ID)
	log.Print("User Name:", updateUser.FirstName, updateUser.LastName)
	log.Print("User Email:", updateUser.Email)
	log.Print("Created at", updateUser.CreatedAt)

}

func TestUsers_WithUsersData(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jason", "last_name":"park", "email":"jason@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
	fmt.Print(m[1])
	fmt.Print(m[0])

}
