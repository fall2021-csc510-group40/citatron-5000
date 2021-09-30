/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package main

import (
	"core/server"
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("p", 80, "port")

	flag.Parse()

	s := server.New()
	panic(s.ListenAndServe(fmt.Sprintf(":%d", *port)))
}
