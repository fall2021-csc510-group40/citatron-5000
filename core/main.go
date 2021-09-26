package main

import (
	"core/server"
)

func main() {
	s := server.New()
	panic(s.ListenAndServe(":8080"))
}
