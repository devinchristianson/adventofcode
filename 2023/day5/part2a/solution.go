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

type Mapping struct {
	start  int // where the mapping takes effect
	end    int // where the mapping ends
	offset int // offset from src -> dst
}
type Range struct {
	start int
	end   int
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
	var initialSeeds []Range
	for scanner.Scan() {
		line := scanner.Text()
		if header := matchMapHeader.FindStringSubmatch(line); len(header) > 0 {
			mapping_src = header[1]
			mapping_dst = header[2]
		} else if initialSeedsLine := matchInitialSeeds.FindAllStringSubmatch(line, -1); len(initialSeedsLine) > 0 {
			values := strings.Fields(strings.Split(line, ":")[1])
			for i := 0; i < len(values)-1; i += 2 {
				seed_start, err := strconv.Atoi(values[i])
				if err != nil {
					panic("bad seed_start number")
				}
				seed_length, err := strconv.Atoi(values[i+1])
				if err != nil {
					panic("bad seed_end number")
				}
				initialSeeds = append(initialSeeds, Range{
					start: seed_start,
					end:   seed_start + seed_length - 1,
				})
			}
		} else if mapLine := matchMapLine.FindAllStringSubmatch(line, -1); len(mapLine) > 0 {
			values := strings.Fields(line)
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
		}
	}
	category := "seed"
	var converter func(category string, value int) int
	converter = func(category string, value int) int {
		//if category != "seed" {
		//	fmt.Printf("to %d\n", value)
		//}
		if category == "location" {
			return value
		}
		//fmt.Printf("Converting %s", category)
		//this for loop will only ever run once, it's just the simplest way to get the next (and only) key in the dict to covert to
		for k, v := range mappings[category] {
			//fmt.Printf("-to-%s: %d ", k, value)
			for _, m := range v {
				if value >= m.start && value <= m.end {
					return converter(k, value+m.offset)
				}
			}
			return converter(k, value)
		}
		panic("hit end of recursive function somehow")
	}
	min_loc := math.MaxInt
	for _, r := range initialSeeds {
		for s := 0; s+r.start < r.end; s++ {
			result := converter(category, s+r.start)
			if result < min_loc {
				min_loc = result
			}
		}
		fmt.Print(".")
	}
	fmt.Println(min_loc)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
