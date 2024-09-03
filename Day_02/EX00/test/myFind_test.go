package test

import (
	"bufio"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestMyFinde1(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myFind.go", "../")
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
		"../cmd",
		"../cmd/myFind.go",
		"../sources",
		"../sources/softlink -> ../sources/source",
		"../sources/softlinkDir -> [broken]",
		"../sources/source",
		"../test",
		"../test/myFind_test.go",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyFinde2(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myFind.go", "-d", "../")
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
		"../cmd",
		"../sources",
		"../test",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyFinde3(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myFind.go", "-f", "-ext", "go", "../")
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
		"../cmd/myFind.go",
		"../test/myFind_test.go",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyFinde4(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myFind.go", "-sl", "../")
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
		"../sources/softlink -> ../sources/source",
		"../sources/softlinkDir -> [broken]",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}

func TestMyFinde5(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	go func() {
		defer w.Close()
		cmd := exec.Command("go", "run", "../cmd/myFind.go", "-d", "-f", "-ext", "go", "-sl", "../../")
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
		"../../EX00",
		"../../EX00/cmd",
		"../../EX00/cmd/myFind.go",
		"../../EX00/sources",
		"../../EX00/sources/softlink -> ../../EX00/sources/source",
		"../../EX00/sources/softlinkDir -> [broken]",
		"../../EX00/test",
		"../../EX00/test/myFind_test.go",
		"../../EX01",
		"../../EX01/cmd",
		"../../EX01/cmd/myWc.go",
		"../../EX01/test",
		"../../EX01/test/myWc_test.go",
		"../../EX01/txt",
		"../../EX02",
		"../../EX02/cmd",
		"../../EX02/cmd/myXargs.go",
		"../../EX02/test",
		"../../EX02/test/myXargs_test.go",
		"../../EX03",
		"../../EX03/archive",
		"../../EX03/cmd",
		"../../EX03/cmd/myRotate.go",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}
