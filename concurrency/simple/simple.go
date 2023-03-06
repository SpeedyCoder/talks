package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		go fmt.Printf("This is job: %v\n", i)
	}
}
