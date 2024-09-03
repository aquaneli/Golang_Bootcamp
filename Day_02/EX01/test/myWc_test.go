package test

import (
	"bufio"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestMyWc1(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myWc.go", "-l", "../txt/test.txt")
		cmd.Stdout = w
		err = cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()

	data := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	expected := []string{
		"3	../txt/test.txt",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyWc2(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myWc.go", "-m", "../txt/test.txt")
		cmd.Stdout = w
		err = cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()

	data := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	expected := []string{
		"751	../txt/test.txt",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyWc3(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myWc.go", "-w", "../txt/test.txt")
		cmd.Stdout = w
		err = cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}()

	data := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	expected := []string{
		"128	../txt/test.txt",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}
