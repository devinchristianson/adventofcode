package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Mapping struct {
	start  int // where the mapping takes effect
	end    int // where the mapping ends
	offset int // offset from src -> dst
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	matchMapHeader := regexp.MustCompile(`^(\w*)-to-(\w*) map:$`)
	matchMapLine := regexp.MustCompile(`^(\d+\s*)+$`)
	matchInitialSeeds := regexp.MustCompile(`^seeds: (\d+\s*)+$`)
	mappings := make(map[string]map[string][]Mapping)
	mapping_src := ""
	mapping_dst := ""
	var initialSeeds []int
	for scanner.Scan() {
		line := scanner.Text()
		if header := matchMapHeader.FindStringSubmatch(line); len(header) > 0 {
			mapping_src = header[1]
			mapping_dst = header[2]
		} else if initialSeedsLine := matchInitialSeeds.FindAllStringSubmatch(line, -1); len(initialSeedsLine) > 0 {
			values := strings.Fields(strings.Split(line, ":")[1])
			initialSeeds = make([]int, len(values))
			for i, v := range values {
				initialSeeds[i], err = strconv.Atoi(v)
				if err != nil {
					panic("bad initial seed number")
				}
			}
		} else if mapLine := matchMapLine.FindAllStringSubmatch(line, -1); len(mapLine) > 0 {
			values := strings.Fields(line)
			fmt.Println(values)
			dst_start, err := strconv.Atoi(values[0])
			if err != nil {
				panic("bad number")
			}
			src_start, err := strconv.Atoi(values[1])
			if err != nil {
				panic("bad number")
			}
			range_length, err := strconv.Atoi(values[2])
			if err != nil {
				panic("bad number")
			}
			if mappings[mapping_src] == nil {
				mappings[mapping_src] = make(map[string][]Mapping)
			}
			mappings[mapping_src][mapping_dst] = append(mappings[mapping_src][mapping_dst], Mapping{
				start:  src_start,
				end:    src_start + range_length,
				offset: dst_start - src_start,
			})
		} else if strings.TrimSpace(line) != "" {
			fmt.Print()
		}
	}
	fmt.Println(initialSeeds)
	fmt.Println(mappings)
	category := "seed"
	var converter func(category string, value int) int
	converter = func(category string, value int) int {
		if category != "seed" {
			fmt.Printf("to %d\n", value)
		}
		if category == "location" {
			return value
		}
		fmt.Printf("Converting %s", category)
		//this for loop will only ever run once, it's just the simplest way to get the next (and only) key in the dict to covert to
		for k, v := range mappings[category] {
			fmt.Printf("-to-%s: %d ", k, value)
			for _, m := range v {
				if value >= m.start && value <= m.end {
					return converter(k, value+m.offset)
				}
			}
			return converter(k, value)
		}
		panic("hit end of recursive function somehow")
	}
	locations := make([]int, len(initialSeeds))
	for i, s := range initialSeeds {
		locations[i] = converter(category, s)
		fmt.Println("-----------------------")
	}
	fmt.Println(locations)
	sort.Ints(locations)
	fmt.Println(locations)
	fmt.Println(locations[0])
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
