package main

import "fmt"

func main() {

	d := map[string]string{
		"name": "<UNK>",
	}

	for k, v := range d {
		fmt.Printf("%s = %s\n", k, v)
	}

	i := 0

	for i < 10 {
		i += 1
		defer fmt.Println(i)
		fmt.Println("循环结束")

	}

}
