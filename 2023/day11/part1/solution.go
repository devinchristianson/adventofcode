package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

type Point struct {
	x int
	y int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var points []*Point
	var emptyY []bool
	scanner.Scan()
	emptyX := make([]bool, len(scanner.Text()))
	max := Point{
		x: len(scanner.Text()),
	}
	for i := range emptyX {
		emptyX[i] = true
	}
	{
		y := 0
		for {
			emptyY = append(emptyY, true)
			x := 0
			for _, r := range scanner.Text() {
				switch r {
				case '#':
					points = append(points, &Point{x: x, y: y})
					emptyX[x] = false
					emptyY[y] = false
				}
				x++
			}

			if !scanner.Scan() {
				max.y = y
				break
			}
			y++
		}
	}
	// generate a list of the empty Y values
	var emptyYVals []int
	emptyYs := 0
	for y, isYEmpty := range emptyY {
		if isYEmpty {
			emptyYVals = append(emptyYVals, y+emptyYs)
			emptyYs++
		}
	}
	// generate a list of the empty X values
	var emptyXVals []int
	emptyXs := 0
	for x, isXEmpty := range emptyX {
		if isXEmpty {
			emptyXVals = append(emptyXVals, x+emptyXs)
			emptyXs++
		}
	}
	max.y += emptyYs
	max.x += emptyXs
	//re-map all the points
	for _, p := range points {
		for _, y := range emptyYVals {
			if p.y > y {
				p.y++
			}
		}
		for _, x := range emptyXVals {
			if p.x > x {
				p.x++
			}
		}
	}
	acc := 0
	for y := 0; y <= max.y; y++ {
		for x := 0; x <= max.x; x++ {
			if acc < len(points) && points[acc].x == x && points[acc].y == y {
				fmt.Print(acc)
				acc++
			} else if slices.Contains(emptyXVals, x) || slices.Contains(emptyYVals, y) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	sum := 0
	for i, p := range points {
		for j := i + 1; j < len(points); j++ {
			sum += int(math.Abs(float64(p.x-points[j].x)) + math.Abs(float64(p.y-points[j].y)))
		}
	}
	fmt.Printf("Part1 sum: %d", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
