package main

import (
	"roob.re/gyndns"
	"os"
	"encoding/json"
	"log"
)

func main() {
	gynFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	params := gyndns.Params{}
	err = json.NewDecoder(gynFile).Decode(&params)
	if err != nil {
		log.Fatalf("Error parsing gyndns.json: %v", err)
	}

	gyndns.New(&params).Run()
}
