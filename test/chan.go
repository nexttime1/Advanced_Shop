package main

import "time"

func main() {
	c := make(chan int, 1)
	go func() {
		<-c
	}()
	go func() {
		<-c
	}()
	time.Sleep(time.Millisecond * 1000)
	close(c)
	time.Sleep(time.Millisecond * 1000)
}
