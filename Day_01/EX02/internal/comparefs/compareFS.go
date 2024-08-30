package comparefs

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// ParseFlag обработывает входные данные
func ParseFlag() (*[2][]byte, error) {
	strOld := flag.String("old", "", "file path old")
	strNew := flag.String("new", "", "file path new")
	flag.Parse()

	if len(*strOld) == 0 || len(*strNew) == 0 {
		return nil, errors.New("file path not specified")
	}

	if len(os.Args) != 5 {
		return nil, errors.New("incorrect number of arguments")
	}

	dataOld, err := readFileAll(strOld)
	if err != nil {
		return nil, err
	}
	dataNew, err := readFileAll(strNew)
	if err != nil {
		return nil, err
	}

	return &[2][]byte{dataOld, dataNew}, nil
}

func readFileAll(str *string) ([]byte, error) {
	b, err := os.ReadFile(*str)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Compare сравнивает 2 текстовых файла
func Compare(res *[2][]byte) {
	resOld := strings.Split(string(res[0]), "\n")
	resNew := strings.Split(string(res[1]), "\n")

	oldCheck := make(map[string]struct{})
	newCheck := make(map[string]struct{})

	for _, v := range resOld {
		oldCheck[v] = struct{}{}
	}

	for _, v := range resNew {
		newCheck[v] = struct{}{}
	}

	indexNew := 0
	indexOld := 0

	for _, v := range resNew[indexNew:] {
		_, ok := oldCheck[v]
		indexNew++
		if !ok {
			fmt.Printf("ADDED %s\n", v)
		}
	}

	for _, v := range resOld[indexOld:] {
		_, ok := newCheck[v]
		indexOld++
		if !ok {
			fmt.Printf("REMOVED %s\n", v)
		}
	}
}
