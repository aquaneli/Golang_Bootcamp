package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	logFilePathIndex, flagA, err := Parsing()
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(flag.Args()[logFilePathIndex:]))
	for i := logFilePathIndex; i < len(flag.Args()); i++ {
		go ProcessingTar(i, *flagA, &wg)
	}
	wg.Wait()
}

// Parsing проверяет log файлы или нет и есть ли флаг a.
// Если флаг активен то архивируем данные по тому пути который указан.
func Parsing() (int, *bool, error) {
	flagA := flag.Bool("a", false, "save dir")
	flag.Parse()
	if len(flag.Args()) == 0 {
		return 0, nil, errors.New("please provide the path to the log file")
	}

	logFilePathIndex := 0
	if *flagA {
		if len(flag.Args()) < 2 {
			return 0, nil, errors.New("few arguments")
		}
		logFilePathIndex = 1
	}

	return logFilePathIndex, flagA, nil
}

// ProcessingTar подготовка данных для создания архива
func ProcessingTar(logFilePathIndex int, flagA bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fileInfo, err := os.Stat(flag.Args()[logFilePathIndex])
	if err != nil {
		log.Fatal(err)
	}
	logFilePath := flag.Args()[logFilePathIndex]
	modTime := fileInfo.ModTime().Unix()
	outputFile := fmt.Sprintf("%s_%d.tar.gz", strings.TrimSuffix(logFilePath, filepath.Ext(logFilePath)), modTime)
	if flagA {
		statInfo, err := os.Stat(flag.Args()[0])
		if err != nil {
			log.Fatal(err)
		}

		if !statInfo.IsDir() {
			log.Fatalf("%s it's not dir", flag.Args()[0])
		}

		outputFile = filepath.Join(flag.Args()[0], filepath.Base(outputFile))
		outputDir := filepath.Dir(outputFile)
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

	}
	CreateArchive(fileInfo, outputFile, logFilePath)
	fmt.Printf("Log file archived to %s\n", outputFile)
}

// CreateArchive создание, сжатие данных и копирование в архив
func CreateArchive(fileInfo fs.FileInfo, outputFile, logFilePath string) {
	// Создание архива
	tarFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Открытие лог-файла
	logFile, err := os.Open(logFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Создание заголовка для файла в архиве
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		log.Fatal(err)
	}

	// Запись заголовка в архив
	err = tarWriter.WriteHeader(header)
	if err != nil {
		log.Fatal(err)
	}

	// Копирование содержимого файла в архив
	_, err = io.Copy(tarWriter, logFile)
	if err != nil {
		log.Fatal(err)
	}
}
