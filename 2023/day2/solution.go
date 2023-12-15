package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Color struct {
	min int
}
type Game struct {
	red   Color
	green Color
	blue  Color
}

func testGame(line string, given_cubes Game) int {
	gameMatcher := regexp.MustCompile(`Game (\d*)`)
	input := strings.Split(line, ":")
	_, err := strconv.Atoi(gameMatcher.FindStringSubmatch(input[0])[1])
	if err != nil {
		panic("oops")
	}
	colorMatcher := regexp.MustCompile(`\s*(\d*) (\w*)\s*`)
	rounds := strings.Split(input[1], ";")
	minimums := Game{}
	for _, round := range rounds {
		colors := strings.Split(round, ",")
		for _, c := range colors {
			result := colorMatcher.FindStringSubmatch(c)
			cube_count, err := strconv.Atoi(result[1])
			cube_color := result[2]
			if err != nil {
				panic("oops")
			}
			switch cube_color {
			case "red":
				if cube_count > minimums.red.min {
					minimums.red.min = cube_count
				}
			case "blue":
				if cube_count > minimums.blue.min {
					minimums.blue.min = cube_count
				}
			case "green":
				if cube_count > minimums.green.min {
					minimums.green.min = cube_count
				}
			}
		}
	}
	return minimums.red.min * minimums.blue.min * minimums.green.min
}
func main() {
	given_cubes := Game{
		red: Color{
			min: 12,
		},
		green: Color{
			min: 13,
		},
		blue: Color{
			min: 14,
		},
	}
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	for scanner.Scan() {
		sum += testGame(scanner.Text(), given_cubes)
	}
	fmt.Printf("The sum is %d\n", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
