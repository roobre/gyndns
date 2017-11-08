package main

import (
	"roob.re/gyndns"
	"os"
	"io/ioutil"
)

func main() {
	usersFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	users, err := ioutil.ReadAll(usersFile)
	if err != nil {
		panic(err)
	}

	server := gyndns.New(nil, users)

	server.Run()
}
