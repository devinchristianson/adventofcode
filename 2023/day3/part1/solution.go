package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type PartNumber struct {
	number   int
	lineNum  int
	startPos int
	endPos   int
}

func checkNumber(part PartNumber, schematic [][]bool) bool {
	// a bit inefficient here as we're also checking if the number's own coords are special
	for _, y := range []int{part.lineNum - 1, part.lineNum, part.lineNum + 1} {
		if y >= 0 && y < len(schematic) {
			for x := part.startPos - 1; x < part.endPos+2; x++ {
				if x >= 0 && x < len(schematic[y]) {
					if schematic[y][x] {
						return true
					}
				}
			}
		}
	}
	return false
}

func convertNumber(buffer bytes.Buffer, x int, y int) PartNumber {
	num, err := strconv.Atoi(buffer.String())
	if err != nil {
		panic("bad number")
	}
	return PartNumber{
		number:   num,
		startPos: x - buffer.Len(),
		endPos:   x - 1,
		lineNum:  y,
	}
}

func main() {
	file, err := os.Open("../input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	schematic := [][]bool{}
	partNumbers := []PartNumber{}
	isDigit := regexp.MustCompile(`^[0-9]$`)
	isSpecialChar := regexp.MustCompile(`[^\d\.]`)
	y := 0
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		convertedLine := make([]bool, len(line))
		for x, c := range line {
			if isDigit.MatchString(string(c)) {
				buffer.WriteRune(c)
			} else {
				if buffer.Len() != 0 {
					partNumbers = append(partNumbers, convertNumber(buffer, x, y))
					buffer.Reset()
				}
			}
			if isSpecialChar.MatchString(string(c)) {
				convertedLine[x] = true
			} else {
				convertedLine[x] = false
			}
		}
		schematic = append(schematic, convertedLine)
		if buffer.Len() != 0 {
			partNumbers = append(partNumbers, convertNumber(buffer, len(line)-1, y))
			buffer.Reset()
		}
		y++
	}
	sum := 0
	for _, p := range partNumbers {
		if checkNumber(p, schematic) {
			sum += p.number
		}

	}
	fmt.Printf("The sum is %d\n", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
