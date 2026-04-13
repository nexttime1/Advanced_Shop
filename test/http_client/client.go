package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	response, err := http.Post("http://localhost:8080/hello", "application/json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	respBody, _ := io.ReadAll(response.Body)
	fmt.Printf("resp: %s", respBody)

}
