package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Lense struct {
	label       string
	focalLength int
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

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	matcher := regexp.MustCompile(`(\w*)([=-])([\d]*)`)
	boxes := map[int][]Lense{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		// categorize all of the lenses
		for _, s := range line {
			groups := matcher.FindStringSubmatch(s)
			label := groups[1]
			hash := hash(label)
			// check if label is already in box
			existingIndex := slices.IndexFunc(boxes[hash], func(l Lense) bool {
				return l.label == label
			})
			switch groups[2] {
			case "=":
				lense := Lense{
					label:       label,
					focalLength: toInt(groups[3]),
				}
				if existingIndex != -1 {
					boxes[hash][existingIndex] = lense
				} else {
					boxes[hash] = append(boxes[hash], lense)
				}
			case "-":
				if existingIndex != -1 {
					boxes[hash] = append(boxes[hash][:existingIndex], boxes[hash][existingIndex+1:]...)
				}

			}
		}
		// calculate focusing power
		for hash, lenses := range boxes {
			for i, lense := range lenses {
				sum += (hash + 1) * (i + 1) * lense.focalLength
			}
		}
	}
	fmt.Println(sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
