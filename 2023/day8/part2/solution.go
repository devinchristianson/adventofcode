package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Node struct {
	right string
	left  string
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers ...int) int {
	if len(integers) < 2 {
		log.Fatal("too few integers")
	}
	result := integers[0] * integers[1] / GCD(integers[0], integers[1])

	for i := 2; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func lastChar(s string) string {
	return string(s[len(s)-1])
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // advance for directions
	directions := scanner.Text()
	scanner.Scan() //advances for whitespace line

	nodeMatcher := regexp.MustCompile(`(\w{3}) = \((\w{3}), (\w{3})\)`)
	nodes := make(map[string]Node)
	var locations []string
	for scanner.Scan() {
		values := nodeMatcher.FindStringSubmatch(scanner.Text())
		nodes[values[1]] = Node{
			left:  values[2],
			right: values[3],
		}
		if lastChar(values[1]) == "A" {
			locations = append(locations, values[1])
		}

	}
	travellerSteps := make([]int, len(locations))
	for i, location := range locations {
		steps := 0
		for lastChar(location) != "Z" {
			for _, d := range directions {
				steps++
				switch string(d) {
				case "L":
					location = nodes[location].left
				case "R":
					location = nodes[location].right
				}
				if lastChar(location) == "Z" {
					travellerSteps[i] = steps
					break
				}
			}
		}
	}

	fmt.Println(LCM(travellerSteps...))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
