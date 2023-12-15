package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	wins   int
	copies int
}

func computeCard(line string, cards map[int]Card) int {
	input := strings.Split(line, ":")
	numbers := strings.Split(input[1], "|")
	cardMatcher := regexp.MustCompile(`Card\s*(\d*)`)
	id, err := strconv.Atoi(cardMatcher.FindStringSubmatch(input[0])[1])
	if err != nil {
		panic("bad number")
	}
	winners := strings.Fields(strings.TrimSpace(numbers[0]))
	drawn := strings.Fields(strings.TrimSpace(numbers[1]))
	matches := 0

	// O(n^2) time but oh well
	for _, d := range drawn {
		for _, w := range winners {
			if w == d {
				matches++
			}
		}
	}
	cards[id] = Card{
		wins:   matches,
		copies: 1,
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

	cards := make(map[int]Card)
	for scanner.Scan() {
		computeCard(scanner.Text(), cards)
	}
	sum := 0
	for k := 1; k <= len(cards); k++ {
		v := cards[k]
		for i := 1; i <= v.wins; i++ {

			if entry, ok := cards[k+i]; ok {
				entry.copies += v.copies
				cards[k+i] = entry
			}
		}
		sum += v.copies
	}
	fmt.Printf("The sum is %d\n", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
