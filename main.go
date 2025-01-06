package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID        string
	Firstname string
	Lastname  string
	Email     string
	Age       uint
	Created   time.Time
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
