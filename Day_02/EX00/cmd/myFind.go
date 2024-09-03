package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Flags структура со статусом флагов
type Flags struct {
	flagD   bool
	flagF   bool
	flagSL  bool
	flagEXT string
}

func main() {
	flags, err := Parse()
	if err != nil {
		log.Fatal(err)
	}
	SearchAndPrint(*flags)
}

// Parse обрабатывает аргументы командной строки
func Parse() (*Flags, error) {
	argD := flag.Bool("d", false, "search only directorys")
	argF := flag.Bool("f", false, "search only files")
	argSL := flag.Bool("sl", false, "search only links")
	argEXT := flag.String("ext", "", "print files with specific extensions")
	flag.Parse()

	if len(flag.Args()) != 1 {
		return nil, errors.New("one path must be specified")
	}

	if !*argF {
		*argEXT = ""
	}

	if !*argD && !*argF && !*argSL {
		*argD, *argF, *argSL = true, true, true
	}

	return &Flags{
		flagD:   *argD,
		flagF:   *argF,
		flagSL:  *argSL,
		flagEXT: *argEXT,
	}, nil
}

// SearchAndPrint рекурсивно проходим по всем папкам от корневой
func SearchAndPrint(flags Flags) {
	filepath.Walk(os.Args[len(os.Args)-1], func(path string, info os.FileInfo, err error) error {
		if os.Args[len(os.Args)-1] == path {
			return nil
		}
		if err != nil {
			return filepath.SkipDir
		}
		if info.IsDir() && flags.flagD {
			fmt.Println(path)
		} else {
			PrintFiles(flags, path, info)
		}
		return nil
	})
}

// PrintFiles печатает файлы или ссылки
func PrintFiles(flags Flags, path string, info os.FileInfo) {
	if info.Mode().Type() == os.ModeSymlink && flags.flagSL {
		link, err := filepath.EvalSymlinks(path)
		if err != nil {
			link = "[broken]"
		}
		fmt.Println(path + " -> " + link)
	} else if flags.flagF {
		if flags.flagEXT != "" {
			split := strings.Split(info.Name(), ".")
			fileFormat := split[len(split)-1]
			if flags.flagEXT == fileFormat {
				fmt.Println(path)
			}
		} else {
			fmt.Println(path)
		}

	}
}
