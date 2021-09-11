package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPage(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL) //ts의 url을 get으로 호출
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))

}

//2차 테스트 : 버퍼에 로거가 잘 찎혔는지 확인. setoutput으로.
func TestDecoHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	buf := &bytes.Buffer{} //아래 한줄씩 버퍼에서 읽도록 만들어서 로거에 남길 수 있도록 하려면?
	log.SetOutput(buf)     //output destination = buf(SetOutput sets the output destination for the standard logger)

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	r := bufio.NewReader(buf)
	line, _, err := r.ReadLine() //한줄씩 읽는당
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Started")

	/*	line, _, err = r.ReadLine()
		assert.NoError(err)
		assert.Contains(string(line), "[LOGGER1] Started")

		line, _, err = r.ReadLine()
		assert.NoError(err)
		assert.Contains(string(line), "[LOGGER1] Completed")

		line, _, err = r.ReadLine()
		assert.NoError(err)
		assert.Contains(string(line), "[LOGGER2] Completed")
	*/
	log.Print(string(line))
}
