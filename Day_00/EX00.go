package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	printData()
}

func printData() {
	buff, num := scanVal()
	if len(buff) > 0 && num >= 1 && num <= 4 {
		sort.Ints(buff)
		arifm := sum(buff)
		for i := 0; i < num; i++ {
			switch i {
			case 0:
				fmt.Printf("%.2f\n", arifm)
			case 1:
				fmt.Printf("%.2f\n", median(buff))
			case 2:
				fmt.Printf("%d\n", mode(buff))
			case 3:
				fmt.Printf("%.2f\n", stanDev(buff, arifm))
			}

		}
	} else {
		log.Fatalln("incorrect value")
	}
}

func scanVal() ([]int, int) {
	scanner := bufio.NewScanner(os.Stdin)
	buff := make([]int, 0, 256)
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) > 0 {
			num, err := strconv.Atoi(str)
			if err != nil || num < -100000 || num > 100000 {
				log.Fatalln("incorrect value")
			}
			buff = append(buff, num)
		} else {
			break
		}
	}

	f := flag.String("f", "4", "some description")
	flag.Parse()

	num, err := strconv.Atoi(*f)
	if err != nil {
		log.Fatalln("incorrect flag")
	}
	return buff, num
}

func sum(buff []int) float64 {
	res := 0
	for _, val := range buff {
		res += val
	}
	return float64(res) / float64(len(buff))
}

func median(buff []int) float64 {
	res := 0
	if len(buff) > 1 {
		res = buff[len(buff)/2-1] + buff[len(buff)/2]
	} else if len(buff) == 1 {
		res = buff[len(buff)/2]
	}
	return float64(res) / 2
}

func mode(buff []int) int {
	res := make(map[int]int, 256)
	max := 0
	minKey := buff[0]
	for _, val := range buff {
		_, ok := res[val]
		if !ok {
			res[val] = 1
		} else {
			res[val]++
			if max < res[val] {
				max = res[val]
				minKey = val
			}
		}
	}
	return minKey
}

func stanDev(buff []int, arifm float64) float64 {
	res := 0.0
	for _, val := range buff {
		res += math.Pow(float64(val)-arifm, 2)
	}
	if len(buff) > 1 {
		res = math.Sqrt(res / float64(len(buff)-1))
	}
	return res
}
