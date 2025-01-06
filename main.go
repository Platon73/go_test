package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Age       uint      `json:"age"`
	Created   time.Time `json:"created"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	users := []*User{}

	rows, err := h.DB.Query("SELECT id, firstname, lastname, email, age, created from users")

	__err_panic(err)

	for rows.Next() {
		post := &User{}
		log.Println("Записи ", rows)
		err = rows.Scan(&post.ID, &post.Firstname, &post.Lastname, &post.Email, &post.Age, &post.Created)
		__err_panic(err)
		users = append(users, post)
	}

	rows.Close()

	var strRes strings.Builder

	if users != nil && len(users) > 0 {
		for _, user := range users {
			strRes.WriteString(user.ID + " " + user.Firstname + " " + user.Lastname + " " + user.Email + " " + user.Created.Format(time.RFC3339) + "\n")
		}
	}

	// Создаем структуру для передачи в шаблон
	data := struct {
		Data string
	}{
		Data: strRes.String(),
	}

	w.Header().Set("Content-Type", "text/html")
	err1 := h.Tmpl.Execute(w, data) // Передаем данные в шаблон
	if err1 != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
}

func getUserForBody(w http.ResponseWriter, r *http.Request) User {
	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
		return User{}
	}
	defer r.Body.Close() // Закрываем тело запроса после чтения

	// Создаем экземпляр структуры User
	var user User

	// Преобразуем JSON в структуру
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Ошибка при преобразовании JSON", http.StatusBadRequest)
		return User{}
	}

	fmt.Printf("Полученный пользователь: %+v\n", user)
	return user
}

func (h *Handler) InsertUser(user User) {
	// в целям упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"INSERT INTO users (id, firstname, lastname, email, age, created) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Age,
		time.Now(),
	)
	__err_panic(err)

	_, err = result.RowsAffected()
	__err_panic(err)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {

	user := getUserForBody(w, r)

	h.InsertUser(user)

	http.Redirect(w, r, "/user", http.StatusFound)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := getUserForBody(w, r)

	// в целям упрощения примера пропущена валидация
	num, err := strconv.Atoi(user.ID)
	__err_panic(err)
	h.DeleteFromId(num)
	h.InsertUser(user)

	http.Redirect(w, r, "/user", http.StatusFound)
}

func (h *Handler) DeleteFromId(id int) {
	fmt.Println("ID для удаления:", id)

	result, err := h.DB.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	__err_panic(err)

	_, err = result.RowsAffected()
	__err_panic(err)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр id из строки запроса
	query := r.URL.Query()
	idStr := query.Get("id") // Получаем значение параметра "id"

	// Преобразуем строку в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	__err_panic(err)

	h.DeleteFromId(id)

	http.Redirect(w, r, "/user", http.StatusFound)
}

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func main() {

	db, err := sql.Open("postgres", "postgres://userok:p@ssw0rd@localhost:5400/pogreb?sslmode=disable")

	if err != nil {
		fmt.Println("db.Prepare failed:", err)
		return
	}

	db.SetMaxOpenConns(10)

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при проверке соединения: %v", err)
	}

	log.Println("Успешно подключено к базе данных!")

	handlers := &Handler{
		DB: db,
		Tmpl: template.Must(template.New("example").Parse("<!DOCTYPE html>" +
			"<html><head>" +
			"<title>List users</title>" +
			"</head>" +
			"<body>" +
			"<h1>List users in JSON</h1>" +
			"{{ .Data}}</body></html>")),
		//Tmpl: template.Must(template.ParseGlob("../crud_templates/*")),
	}

	// в целям упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/user", handlers.List).Methods("GET")
	r.HandleFunc("/user/create", handlers.Add).Methods("POST")
	r.HandleFunc("/user/update", handlers.Update).Methods("PATCH")
	r.HandleFunc("/user/delete", handlers.Delete).Methods("DELETE")

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)
}

// не используйте такой код в прошакшене
// ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
