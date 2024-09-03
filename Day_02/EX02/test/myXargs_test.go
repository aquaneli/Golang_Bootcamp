package test

import (
	"bufio"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestMyXargs1(t *testing.T) {
	r1, w1, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r1.Close()

	go func() {
		defer w1.Close()

		cmd1 := exec.Command("./myFind", "-f", "-ext", "go", "../../")
		cmd1.Stdout = w1

		err = cmd1.Run()
		if err != nil {
			t.Error(err)
		}
	}()

	r2, w2, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	defer r2.Close()
	go func() {
		defer w2.Close()
		cmd2 := exec.Command("go", "run", "../cmd/myXargs.go", "./myWc")
		cmd2.Stdout = w2
		cmd2.Stdin = r1
		err = cmd2.Run()
		if err != nil {
			t.Error(err)
		}
	}()

	data := []string{}
	scanner := bufio.NewScanner(r2)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	expected := []string{
		"264	../../EX00/cmd/myFind.go",
		"440	../../EX00/test/myFind_test.go",
		"242	../../EX01/cmd/myWc.go",
		"241	../../EX01/test/myWc_test.go",
		"143	../../EX02/cmd/myXargs.go",
		"141	../../EX02/test/myXargs_test.go",
		"335	../../EX03/cmd/myRotate.go",
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %s but get %s", expected, data)
	}
}
