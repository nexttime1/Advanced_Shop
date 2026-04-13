package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("成功响应")

		w.Write([]byte("Hello World!"))

	})

	http.ListenAndServe(":8080", nil)

}
