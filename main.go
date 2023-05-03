package main

import (
	// "fmt"
)

func main() {
	// fmt.Println("alhamdulillah")

	server := NewAPIServer(":3000")
	server.Run()
}