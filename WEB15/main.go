package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback", // 로그인 완료 후의 리다이렉트. 구글에서 알려준 결과, callback을 알려줄 곳. URL to redirect users going through the OAuth flow, after the resource owner's URLs.
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"}, // 이메일 권한 요청. 요청하는 데이터의 스코프. 권한 범위=auth=authority의 이메일.
	Endpoint:     google.Endpoint,
}

// 111. googleLoginHandler 호출 -> 유저가 로그인 페이지에 접속
// 222. 로그인 페이지 접속 시 유저를 식별하기 위해 생성한 랜덤한 state값을 사용해 구글 로그인 링크를 생성한다.
func googleLoginHandler(w http.ResponseWriter, r *http.Request) { // 구글 로그인 요청을 받으면 googleOauthCoonfig를 통해 구글의 어떤경로로 보내야 되는 지가 나오고, 유저가 그 경로로 접근해서 로그인할 수 있도록 리다이렉트
	state := generateStateOauthCookie(w)                   // AuthCodeURL에 들어가는 state.
	url := googleOauthConfig.AuthCodeURL(state)            // state 받아서 로그인을 위한 URL반환. oath 제공자의 동의페이지(요구된 스코프에 대한 허가) url 반환.
	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 333. 이 경로로 리다이렉트 -> 거기서 로그인폼이 뜬다 => 완료되면 CallBack URL을 구글에서 알려준다(내가 구글 클라우드에서 설정해둔 주소)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) // 현재 + 24시간

	b := make([]byte, 16) //16바이트 array
	rand.Read(b)          //
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration} //state가 들어갈 cookie struct 생성
	http.SetCookie(w, cookie)                                                     // write에 cookie 넣어준다
	return state
}

// 444. CallBack URL로 리다이렉트 => googleAuthCallback 핸들러 호출
// 555. state값이 이전 값과 같은지 확인.
func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate") // 쿠키 읽어오고

	if r.FormValue("state") != oauthstate.Value { // request에 state 넣어주기 위해 FormValue 사용.
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state")) //공격시도가 있었다고 판단되므로 log 남김
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)                                                // 에러 반환하면 해커에게 단서를 제공하게 되므로 그냥 절대경로로 리다이렉스 시켜버림
		return
	}
	//state값이 이전 값과 같으면 유저정보 요청
	data, err := getGoogleUserInfo(r.FormValue("code")) // 여기에서 구글이 code를 알려줌
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data))
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code) // code를 받아서 token으로 교환
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken) // 유저정보를 리퀘스트 하는 경로 "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	if err != nil {
		return nil, fmt.Errorf("Dailed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body) // 위의 유저정보 리퀘스트 경로의 바디값, 즉 유저정보를 불러옴.
}

func main() {
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
