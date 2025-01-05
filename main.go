package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
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

	rows, err := h.DB.Query("SELECT * from user")

	__err_panic(err)

	for rows.Next() {
		post := &User{}
		err = rows.Scan(&post.ID, &post.Firstname, &post.Lastname, &post.Email, &post.Age)
		__err_panic(err)
		users = append(users, post)
	}

	rows.Close()

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalf("Ошибка при маршализации пользователей: %v", err)
	}

	w.Write(usersJSON)
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
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("../user/*")),
	}

	// в целям упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.List).Methods("GET")
	//r.HandleFunc("/items", handlers.List).Methods("GET")
	//r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	//r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	//r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	//r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	//r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

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
