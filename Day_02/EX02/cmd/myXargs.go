package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	err := ParseCommand()
	if err != nil {
		log.Fatal(err)
	}
	err = Processing()
	if err != nil {
		log.Fatal(err)
	}
}

// ParseCommand парсит аргументы myXargs как отдельную команду
func ParseCommand() error {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return errors.New("few arguments")
	}
	return nil
}

// Processing из стандартного потока ввода берем данные и обрабатываем в нашем запущенном процессе
func Processing() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		args := append(flag.Args()[1:], input)
		cmd := exec.Command(flag.Args()[0], args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
