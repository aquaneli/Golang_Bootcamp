package readdb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// DBReader интерфейс для публичного взаимодействия
type DBReader interface {
	dataDeserialization([]byte) (*recipes, error)
	dataSerializationAnotherFormat(recipes) ([]byte, error)
}

type recipes struct {
	Cake []cake `xml:"cake"    json:"cake"`
}

type cake struct {
	Name        string        `xml:"name"    json:"name"`
	Stovetime   string        `xml:"stovetime"    json:"time"`
	Ingredients []ingredients `xml:"ingredients>item"    json:"ingredients"`
}

type ingredients struct {
	Itemname  string `xml:"itemname"    json:"ingredient_name"`
	Itemcount string `xml:"itemcount"    json:"ingredient_count"`
	Itemunit  string `xml:"itemunit"    json:"ingredient_unit"`
}

type dbReaderXML struct{}

func (dbReaderXML) dataDeserialization(b []byte) (*recipes, error) {
	r := recipes{}
	err := xml.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (dbReaderXML) dataSerializationAnotherFormat(r recipes) ([]byte, error) {
	b, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return nil, err
	}
	return b, nil
}

type dbReaderJSON struct{}

func (dbReaderJSON) dataDeserialization(b []byte) (*recipes, error) {
	r := recipes{}
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (dbReaderJSON) dataSerializationAnotherFormat(r recipes) ([]byte, error) {
	b, err := xml.MarshalIndent(r, "", "    ")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Deserialization  публичная функция для взаимодействия, с помощью
// которой десериализуется из формата JSON или XML в структуру
func Deserialization(r DBReader, b []byte) (interface{}, error) {
	obj, err := r.dataDeserialization(b)
	if err != nil {
		return nil, err
	}
	return *obj, nil
}

// SerializationAnotherFormat  публичная функция для взаимодействия, с помощью
// которой сериализуется из формата JSON или XML в структуру
func SerializationAnotherFormat(reader DBReader, r interface{}) ([]byte, error) {
	if reflect.TypeOf(r) != reflect.TypeFor[recipes]() {
		return nil, errors.New("incorrect type")
	}
	b, err := reader.dataSerializationAnotherFormat(r.(recipes))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ParseFlag обработывает входные данные
func ParseFlag() (DBReader, []byte, error) {
	str := flag.String("f", "", "file path")
	flag.Parse()

	if len(*str) == 0 {
		return nil, nil, errors.New("file path not specified")
	}

	if len(os.Args) != 3 {
		return nil, nil, errors.New("incorrect number of arguments")
	}

	format := strings.Split(*str, ".")
	dataFromFile, err := readFileAll(str)
	if err != nil {
		return nil, nil, err
	}

	if format[len(format)-1] == "xml" {
		return dbReaderXML{}, dataFromFile, nil
	} else if format[len(format)-1] == "json" {
		return dbReaderJSON{}, dataFromFile, nil
	}

	return nil, nil, errors.New("incorrect format")
}

func readFileAll(str *string) ([]byte, error) {
	b, err := os.ReadFile(*str)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// PrintData вывод в Stdout в формате JSON или XML
func PrintData(data []byte) {
	str := string(data)
	for _, v := range str {
		fmt.Print(string(v))
	}
}
