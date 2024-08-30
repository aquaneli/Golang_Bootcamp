package main

import (
	"EX00/internal/readdb"
	"log"
)

func main() {
	arg, data, err := readdb.ParseFlag()
	if err != nil {
		log.Fatalln(err)
	}
	resRecipes, err := readdb.Deserialization(arg, data)
	if err != nil {
		log.Fatalln(err)
	}
	resBytes, err := readdb.SerializationAnotherFormat(arg, resRecipes)
	if err != nil {
		log.Fatalln(err)
	}
	readdb.PrintData(resBytes)
}
