package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request) { //파일 업로드 핸들러 정의
	//상대방이 전송해준 파일을 리딩//
	uploadFile, header, err := r.FormFile("upload_file") //들어오는 값이 키 값. 키를 넣으니까 벨류가 나옴. 데이터를 임시로 저장(업로드파일, 헤더에) 있는 거시 폼파일
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close() //핸들러 자원을 사용 후 defer로 닫아 줘야 함.os의 자원을 가져왔기 때문에(open) 닫아줘야 함여.
	//defer는 close가 반환될때 까지 실행 지연(????? 알아보자)

	dirname := "./uploads"
	os.MkdirAll(dirname, 0777)                                 //8진수(리눅스 명령어) - 의미? read write all! 모두 할수있다!
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename) //폼파일 - %s디렉토리 경로/%s파일이름
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	//파일을 카피하는 공간
	io.Copy(file, uploadFile) //위에서 업로드 된 파일(file)을 아래 비어 있는 경로(uploadFile)에 카피
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath) //filepath(sprintf에서 가져온 값)(%s%s)는 upload 경로를 알려준다
}

func main() {
	http.HandleFunc("/uploads", uploadsHandler)           //파일 업로드 핸들러를 만들어준다. 목적지가 업로드 핸들러
	http.Handle("/", http.FileServer(http.Dir("public"))) //목적지가 파일서버의 퍼블릭 디렉토리

	http.ListenAndServe(":3000", nil)
}
