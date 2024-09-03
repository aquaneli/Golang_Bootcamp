package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

// Flags структура со статусом флагов
type Flags struct {
	flagL      bool
	flagM      bool
	flagW      bool
	filesCount int
}

func main() {
	flags, err := Parse()
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(flags.filesCount)

	for i := 0; i < len(flag.Args()); i++ {
		go ParseFile(&wg, *flags, i)
	}
	wg.Wait()
}

// ParseFile обрабатывает файлы
func ParseFile(wg *sync.WaitGroup, flags Flags, index int) {
	b, err := os.ReadFile(flag.Args()[index])
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	count := 0
	reader := strings.NewReader(string(b))
	scanner := bufio.NewScanner(reader)

	if flags.flagL {
		for scanner.Scan() {
			count++
		}
	} else if flags.flagM {
		count = utf8.RuneCountInString(string(b))
	} else if flags.flagW {
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			count++
		}
	}

	fmt.Printf("%d\t%s\n", count, flag.Args()[index])
	wg.Done()

}

// Parse обрабатывает аргументы командной строки
func Parse() (*Flags, error) {
	argL := flag.Bool("l", false, "quantity strings")
	argM := flag.Bool("m", false, "quantity runes")
	argW := flag.Bool("w", false, "quantity words")
	flag.Parse()

	err := CheckFlags(argW)
	if err != nil {
		return nil, err
	}

	files := len(flag.Args())
	if files == 0 {
		return nil, errors.New("few arguments")
	}

	return &Flags{
		flagL:      *argL,
		flagM:      *argM,
		flagW:      *argW,
		filesCount: files,
	}, nil
}

// CheckFlags проверяет флаги
func CheckFlags(argW *bool) error {
	count := 0
	flag.Visit(func(f *flag.Flag) {
		count++
	})
	if count > 1 {
		return errors.New("incorrect number of arguments")
	}
	if count == 0 {
		*argW = true
	}
	return nil
}
