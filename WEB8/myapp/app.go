package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// User struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userMap map[int]*User
var lastID int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id] //id에 해당하는 유저맵 정보를 받아옴. user에 맵데이터 담감.
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// Created User
	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 //string형 id를 추출해서 map에 넣어준다. 리퀘스트(r) 받아서 맵스트링으로 보냄
	id, err := strconv.Atoi(vars["id"]) //vars(맵 스트링)는 string 이므로 atoi로 int형으로 변환. 반환값은 int타입 id랑 err
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	_, ok := userMap[id] //해당 id의 유저맵 정보를 불러오는데, 일단 무서. ok값에도 넣어
	if !ok {             //이게 !ok, 즉 false인 경우. 맵에 id가 없을 경우(지울게 없을 경우)
		w.WriteHeader(http.StatusOK) //원하는 작업(delete)이 잘 됐다는 소리니까 statusok
		fmt.Fprint(w, "No User ID:", id)
		return
	}
	delete(userMap, id) //위의 걍우가 아니면. 유저맵에 id가 존재하면delete에 맵을 넣고 키를 넣으면 맵에서 키에 해당하는 값이 지워짐. 이 경우는 id
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User ID:", id)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {

	updateUser := new(User) //핸들러를 작동시키면 json데이터가 오도록 위에 user struct를 만들어준다
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[updateUser.ID] //usermap에서 id를 찾아서 usermap에 id가 있는경우, user가 *User이므로(ex updateUser id가 11이면 유저맵의 id에 11을 넣어주고, 그게 있는지 확인)
	if !ok {                           //usermap에 update id가 없는 경우, 즉 false
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", updateUser.ID)
		return //탈출!
	}
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	userMap[updateUser.ID] = user
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	return mux
}
