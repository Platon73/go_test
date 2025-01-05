package main

import (
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Привет мир!!!"))
}

func main() {

}

// не используйте такой код в прошакшене
// ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
