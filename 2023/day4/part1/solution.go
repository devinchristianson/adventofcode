package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func computeCard(line string) int {
	input := strings.Split(strings.Split(line, ":")[1], "|")
	winners := strings.Fields(strings.TrimSpace(input[0]))
	drawn := strings.Fields(strings.TrimSpace(input[1]))
	matches := 0

	// O(n^2) time but oh well
	for _, d := range drawn {
		for _, w := range winners {
			if w == d {
				matches++
			}
		}
	}
	return int(math.Pow(2, float64(matches-1)))
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
		sum += computeCard(scanner.Text())
	}
	fmt.Printf("The sum is %d\n", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
