package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tanimutomo/sqlfile"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "postgres://userok:p@ssw0rd@localhost:5400/pogreb?sslmode=disable")

	db.SetMaxOpenConns(10)

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при проверке соединения: %v", err)
	}

	log.Println("Успешно подключено к базе данных!")

	// Инициализация sqlfile
	s := sqlfile.New()

	// Загрузка SQL-файла
	err = s.File("user.sql")
	if err != nil {
		log.Fatalf("Ошибка загрузки SQL-файла: %v", err)
	}

	// Выполнение загруженных запросов
	_, err = s.Exec(db)
	if err != nil {
		log.Fatalf("Ошибка выполнения SQL-запросов: %v", err)
	}
	log.Println("SQL-скрипт успешно выполнен!")
}
