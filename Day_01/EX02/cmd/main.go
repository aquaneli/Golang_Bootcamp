package main

import (
	"EX02/internal/comparefs"
	"log"
)

func main() {
	res, err := comparefs.ParseFlag()
	if err != nil {
		log.Fatal(err)
	}
	comparefs.Compare(res)
}
