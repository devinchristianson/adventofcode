package main

//for part 2 I just edited the input and ran the program again cuz it seemed easier than editing the code to the parser

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	time     int
	distance int
}

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

func findLimits(time int, distance int) (int, int) {
	i := math.Sqrt(float64(time*time - 4*distance))
	return int(math.Floor((float64(time) - i) / 2.0)), int(math.Ceil((float64(time) + i) / 2.0))
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	times, err := stringsToIntegers(strings.Fields(strings.Split(scanner.Text(), ":")[1]))
	if err != nil {
		panic("bad time number")
	}
	scanner.Scan()
	distances, err := stringsToIntegers(strings.Fields(strings.Split(scanner.Text(), ":")[1]))
	if err != nil {
		panic("bad distance number")
	}
	tally := 1
	if len(times) == len(distances) {
		races := make([]Race, len(distances))
		for i, t := range times {
			races[i] = Race{
				time:     t,
				distance: distances[i],
			}
		}
		for _, race := range races {
			min, max := findLimits(race.time, race.distance)
			tally *= (max - min - 1)
		}
	}
	fmt.Println(tally)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
