package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Point struct {
	x int64
	y int64
}

const expansionFactor = 1000000

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
		x: int64(len(scanner.Text())),
	}
	for i := range emptyX {
		emptyX[i] = true
	}
	{
		y := int64(0)
		for {
			emptyY = append(emptyY, true)
			x := int64(0)
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
	var emptyYVals []int64
	emptyYs := int64(0)
	for y, isYEmpty := range emptyY {
		if isYEmpty {
			emptyYVals = append(emptyYVals, int64(int64(y)+emptyYs))
			emptyYs += (expansionFactor - 1)
		}
	}
	// generate a list of the empty X values
	var emptyXVals []int64
	emptyXs := int64(0)
	for x, isXEmpty := range emptyX {
		if isXEmpty {
			emptyXVals = append(emptyXVals, int64(int64(x)+emptyXs))
			emptyXs += (expansionFactor - 1)
		}
	}
	max.y += emptyYs
	max.x += emptyXs
	// re-map all the points
	for _, p := range points {
		for _, y := range emptyYVals {
			if p.y > y {
				p.y += (expansionFactor - 1)
			}
		}
		for _, x := range emptyXVals {
			if p.x > x {
				p.x += (expansionFactor - 1)
			}
		}
	}
	fmt.Println()
	sum := 0
	for i, p := range points {
		for j := i + 1; j < len(points); j++ {
			sum += int(math.Abs(float64(p.x-points[j].x)) + math.Abs(float64(p.y-points[j].y)))
		}
	}
	fmt.Printf("Part2 sum: %d", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
