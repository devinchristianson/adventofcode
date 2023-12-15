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
	position rune
	number   int
	onGroup  int
}

type StateMachine struct {
	states map[State]int
	next   map[State]int
}

func opposite(s rune) rune {
	switch s {
	case WORKING:
		return BROKEN
	case BROKEN:
		return WORKING
	}
	panic("opposite function has been misused - illegal inpput")
}

func (m *StateMachine) Move(s State, p rune, c int) {
	m.Add(State{
		position: p,
		number:   s.number,
		onGroup:  s.onGroup,
	},
		c)
}

func (m *StateMachine) Add(s State, c int) {
	if m.next[s] > 0 {
		m.next[s] = m.next[s] + c
	} else {
		m.next[s] = c
	}
}

func (m *StateMachine) Advance() {
	m.states, m.next = m.next, m.states // swap
	for s := range m.next {
		delete(m.next, s)
	}
}

func (s State) String() string {
	return fmt.Sprintf("{%s, %d, %d}", string(s.position), s.number, s.onGroup)
}

func runStateMachine(springs string, groups []int) int {
	machine := StateMachine{
		states: map[State]int{{WORKING, 0, 0}: 1},
		next:   map[State]int{},
	}
	permutations := 0
	// adding a cheeky extra working pipe to perform the cleanup for me
	for _, spring := range springs + ".." {
		for state, count := range machine.states {
			switch state.position {
			case BROKEN:
				state.number++
				if state.onGroup == len(groups) || state.number > groups[state.onGroup] {
					continue
				}
			case WORKING:
				if state.onGroup < len(groups) && state.number == groups[state.onGroup] {
					state.onGroup++
				} else if state.number > 0 {
					continue
				}
				state.number = 0
			}
			if spring == state.position || spring == UNKNOWN {
				machine.Add(state, count)
			}
			if spring != state.position || spring == UNKNOWN {
				machine.Move(state, opposite(state.position), count)
			}
		}
		machine.Advance()
	}
	for p, c := range machine.states {
		if p.onGroup == len(groups) {
			permutations += c
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
	part2sum := 0
	part2 := true
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		groups := Map(strings.Split(line[1], ","), toInt)
		springs := line[0]
		part1sum += runStateMachine(springs, groups)
		if part2 {
			var part2Groups []int
			part2Springs := ""
			for i := 0; i < 5; i++ {
				part2Groups = append(part2Groups, groups...)
				part2Springs = part2Springs + "?" + springs
			}
			part2sum += runStateMachine(part2Springs[1:], part2Groups)
		}
	}
	fmt.Println(part1sum)
	fmt.Println(part2sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
