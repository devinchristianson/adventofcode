package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func stringsToIntegers(lines []string) ([]int, error) {
	integers := make([]int, len(lines))
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		integers[i] = n
	}
	return integers, nil
}

func findDifferences(values []int) ([]int, bool, int) {
	deltas := make([]int, len(values)-1)
	zeroes := true
	for i := 0; i < len(values)-1; i++ {
		deltas[i] = values[i+1] - values[i]
		zeroes = zeroes && deltas[i] == 0
	}
	return deltas, zeroes, deltas[len(deltas)-1]
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		values, err := stringsToIntegers(strings.Fields(scanner.Text()))
		if err != nil {
			log.Fatalln("Bad integer found")
		}
		lasts := []int{values[len(values)-1]}
		for values, zeroes, last := findDifferences(values); !zeroes; {
			lasts = append(lasts, last)
			values, zeroes, last = findDifferences(values)
		}
		acc := 0
		for i := len(lasts) - 1; i >= 0; i-- {
			acc += lasts[i]
		}
		sum += acc
	}
	fmt.Println(sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
