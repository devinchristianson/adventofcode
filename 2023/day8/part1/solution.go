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
	for scanner.Scan() {
		values := nodeMatcher.FindStringSubmatch(scanner.Text())
		nodes[values[1]] = Node{
			left:  values[2],
			right: values[3],
		}

	}
	location := "AAA"
	steps := 0
	fmt.Println(directions)
	for location != "ZZZ" {
		for _, d := range directions {
			steps++
			switch string(d) {
			case "L":
				location = nodes[location].left
			case "R":
				location = nodes[location].right
			}
			if location == "ZZZ" {
				break
			}
		}
	}
	fmt.Println(steps)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
