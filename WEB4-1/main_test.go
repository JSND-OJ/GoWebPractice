package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	//1단계 파일 읽고 업로딩 삭제 복구여부
	path := "C:/Users/LeeSK/go/src/GoBlock/goweb/WEB4-1/uploads/aaa.jpg"
	file, _ := os.Open(path) //_언더바는 에러무시
	defer file.Close()

	os.RemoveAll("./uploads") //이전 업로드 저장된 파일을 모두 지운다

	//Form파일 만들기
	buf := &bytes.Buffer{}                                                  //버퍼 인스턴스 바이트값 주소 //데이터가 여기서 대기//버퍼는 작은 저장소(주소값이 들어감)
	writer := multipart.NewWriter(buf)                                      //mime 웹기반 메이에서 파일전송 표준포맷 mulipart(형식에 구애받지 않음)는 파일전송 기본포맷방식
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path)) //Base 끝에 있는 것(aaa.jpg)만) - 파일이름만 잘라내기. filepath.Base지정된 경로의 마지막 요소를 반환
	assert.NoError(err)                                                     //검사, 에러가 없어야 한다. 왜? 복구하기위해
	io.Copy(multi, file)                                                    //destination 은 multi, 어디서 가져오나 file// 여기서 복구
	writer.Close()                                                          //데이터는 전부 buf에 실려있어

	res := httptest.NewRecorder()                                //찐서버면 리코더가 아니라 리스폰스
	req := httptest.NewRequest("POST", "/uploads", buf)          //POST는 CREATED로 반환
	req.Header.Set("Content-type", writer.FormDataContentType()) //inputfile로 전송되는 FormData Content-type 알려주어서 해결.
	//버퍼에 업로드 된 파일의 형식을 알랴줌

	uploadsHandler(res, req)              //콜, 데이터 전송
	assert.Equal(http.StatusOK, res.Code) //검사단계
	data, _ := ioutil.ReadAll(res.Body)   //res.Body - main.go 업로드 핸들러의 결과값, filepath
	log.Print(string(data))

	//2단계 원본파일과 업로드파일 성공여부 검증 테스팅
	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, err = os.Stat(uploadFilePath) //stat은 파일정보를 알려준다
	assert.NoError(err)              //통과되면 파일이 존재한다는 의미

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{} //업로드 데이터 읽어온다
	originData := []byte{} // 오리진 데이터 읽어온다
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(originData, uploadData)

}
