package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func toInt(c string) int {
	i, e := strconv.Atoi(string(c))
	if e != nil {
		panic("bad int")
	}
	return i
}

const (
	UNKNOWN = '?'
	WORKING = '.'
	BROKEN  = '#'
	START   = '>'
)

type State struct {
	position    rune
	number      int
	onGroup     int
	accumulator string
}

type StateMachine struct {
	states map[State]bool
	next   map[State]bool
}

func (m *StateMachine) Move(s State, p rune) {
	m.Add(State{position: p, number: s.number, onGroup: s.onGroup, accumulator: s.accumulator})
}

func (m *StateMachine) Add(s State) {
	m.next[s] = false
}

func (m *StateMachine) Advance() {
	m.states = m.next
	m.next = make(map[State]bool)
}

func (s State) String() string {
	return fmt.Sprintf("{%s, %d, %d, %s}", string(s.position), s.number, s.onGroup, s.accumulator)
}

func runStateMachine(springs string, groups []int) int {
	machine := StateMachine{states: map[State]bool{{WORKING, 0, 0, ""}: false}, next: map[State]bool{}}
	permutations := 0
	// adding a cheeky extra working pipe to perform the cleanup for me
	for _, spring := range springs + "." {
		for state := range machine.states {
			state.accumulator = state.accumulator + string(state.position)
			switch state.position {
			case BROKEN:
				state.number++
				if state.onGroup == len(groups) || state.number > groups[state.onGroup] {
					continue
				} else if state.number == groups[state.onGroup] && spring != BROKEN {
					state.position = WORKING
					state.number = 0
					state.onGroup++
				}
				switch spring {
				case BROKEN:
					machine.Add(state)
				case UNKNOWN:
					// keep one state here w/ one more match
					machine.Add(state)
					// fallthrough to also treat as a non-match
					fallthrough
				case WORKING:
					machine.Move(state, WORKING)
				}
			case WORKING:
				if state.number > 0 {
					continue
				}
				state.number = 0
				switch spring {
				case WORKING:
					machine.Add(state)
				case UNKNOWN:
					machine.Add(state)
					fallthrough
				case BROKEN:
					machine.Move(state, BROKEN)
				}
			}
		}
		machine.Advance()
	}
	for p := range machine.states {
		if p.onGroup == len(groups) {
			permutations++
		}
	}
	return permutations
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part1sum := 0
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		groups := Map(strings.Split(line[1], ","), toInt)
		springs := line[0]
		part1sum += runStateMachine(springs, groups)
	}
	fmt.Println(part1sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
