package main

import (
	"EX01/internal/comparedb"
	"log"
)

func main() {
	args, data, err := comparedb.ParseFlag()
	if err != nil {
		log.Fatalln(err)
	}

	resRecipes, err := comparedb.Deserialization(*args, *data)
	if err != nil {
		log.Fatalln(err)
	}

	comparedb.Comparison(resRecipes)
}
