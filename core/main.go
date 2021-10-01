/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package main

import (
	"core/server"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
)

func readServerConfig(path string) (config *server.Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &config)
	return
}

func main() {
	port := flag.Int("p", 80, "port")
	configPath := flag.String("c", "", "config")

	flag.Parse()

	config, err := readServerConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.New(config)
	if err != nil {
		log.Fatal(err)
	}

	panic(s.ListenAndServe(fmt.Sprintf(":%d", *port)))
}
