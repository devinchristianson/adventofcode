package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Point struct {
	y int
	x int
}

func hash(s string) int {
	acc := 0
	for _, r := range s {
		acc += int(r)
		acc *= 17
		acc %= 256
	}
	return acc
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
		line := strings.Split(scanner.Text(), ",")
		for _, s := range line {
			sum += hash(s)
		}
	}
	fmt.Println(sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
