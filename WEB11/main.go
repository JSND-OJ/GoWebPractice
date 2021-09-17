package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 30
}

func main() {
	user := User{Name: "me", Email: "me@naver.com", Age: 18}
	user2 := User{Name: "bbb", Email: "bbb@gmail.com", Age: 40}
	users := []User{user, user2}
	tmpl, err := template.New("Tmpl1").ParseFiles("templates/tmpl1.tmpl", "templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl1.tmpl", user)  //Execute(첫번째 결과를 어디에 쓸거야(Write=os.Strout=화면에 출력하기), user
	//tmpl.ExecuteTemplate(os.Stdout, "tmpl1.tmpl", user2) // tmpl은 틀은 같고 내용만 달라진다
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
