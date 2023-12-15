package main


import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Variety int

const (
	high_card Variety = iota
	one_pair
	two_pair
	three_of_a_kind
	full_house
	four_of_a_kind
	five_of_a_kind
)

type Hand struct {
	hand    string
	variety Variety
	bid     int
}

func calculateVariety(hand string) Variety {
	jokers := 0
	hand_map := make(map[string]int)
	for _, c := range hand {
		if string(c) == "J" {
			jokers++
		} else {
			hand_map[string(c)]++
		}
	}
	keys := make([]string, 0, len(hand_map))

	for key := range hand_map {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return hand_map[keys[i]] > hand_map[keys[j]]
	})
	// for now just going to assume it's always best to turn all the jokers
	// into whatever we have the most of
	if jokers > 0 {
		if len(keys) > 0 {
			hand_map[keys[0]] += jokers
		} else { //all jokers
			return five_of_a_kind
		}
	}
	if hand_map[keys[0]] == 5 {
		return five_of_a_kind
	} else if hand_map[keys[0]] == 4 {
		return four_of_a_kind
	} else if hand_map[keys[0]] == 3 {
		if len(keys) > 1 && hand_map[keys[1]] == 2 {
			return full_house
		} else {
			return three_of_a_kind
		}
	} else if hand_map[keys[0]] == 2 {
		if len(keys) > 1 && hand_map[keys[1]] == 2 {
			return two_pair
		} else {
			return one_pair
		}
	} else {
		return high_card
	}
}

func cardValue(card rune) int {
	switch string(card) {
	case "A":
		return 14
	case "K":
		return 13
	case "Q":
		return 12
	case "J":
		return 0
	case "T":
		return 10
	default:
		value, err := strconv.Atoi(string(card))
		if err != nil {
			panic("bad number")
		} else {
			return value
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hands []Hand
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), " ")
		bid, err := strconv.Atoi(values[1])
		if err != nil {
			panic("bad bid")
		}
		hands = append(hands, Hand{
			hand:    values[0],
			bid:     bid,
			variety: calculateVariety(values[0]),
		})
	}
	sort.SliceStable(hands, func(i, j int) bool {
		if hands[i].variety != hands[j].variety {
			return hands[i].variety < hands[j].variety
		} else {
			for index, icard := range hands[i].hand {

				if icard != rune(hands[j].hand[index]) {
					return cardValue(icard) < cardValue(rune(hands[j].hand[index]))
				}
			}
			return false // hands are identical, so not less than
		}
	})
	sum := 0
	for i, h := range hands {
		sum += h.bid * (i + 1)
	}
	fmt.Println(sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
