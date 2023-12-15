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

type Range struct {
	start int
	end   int
}

// from https://freshman.tech/snippets/go/concatenate-slices/
// because I didn't want to use cause a bunch of chained append calls
func concatMultipleSlices[T any](slices [][]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
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
	var converter func(category string, value_range Range) []Range
	converter = func(category string, value_range Range) []Range {
		if category == "location" {
			return []Range{value_range}
		}
		var ranges []Range
		var re_process Range
		//this for loop will only ever run once per call, it's just the simplest way to get the next (and only) key in the dict to covert to
		for k, v := range mappings[category] {
			for _, m := range v {
				if value_range.start >= m.start { // if range start is >= mapping start
					if value_range.end <= m.end { // and range end is <= mapping end
						ranges = []Range{{ // the mapping applies to the whole range
							start: value_range.start + m.offset,
							end:   value_range.end + m.offset,
						}}
					} else if value_range.start < m.end { // otherwise, if only the range start is <= mapping end
						ranges = []Range{{ // convert the left half of the range
							start: value_range.start + m.offset,
							end:   m.end - 1 + m.offset,
						}}
						re_process = Range{ // re-process the right half of the range
							start: m.end,
							end:   value_range.end,
						}
					}
				} else if value_range.end <= m.start && value_range.end < m.end { // otherwise, if right only half overlaps with mapping
					ranges = []Range{{
						start: value_range.start,
						end:   m.start - 1,
					}, {
						start: m.start + m.offset,
						end:   value_range.end + m.offset,
					},
					}
				} else if value_range.start < m.start && value_range.end > m.end { // if range fully encaptures mapping
					ranges = []Range{{ // convert the overlapping portion of the range (b/c sorted, no need to re-process)
						start: value_range.start,
						end:   m.start - 1,
					}, { // convert w/ offset the overlapping portion of the range
						start: m.start + m.offset,
						end:   m.end - 1 + m.offset,
					},
					}
					re_process = Range{ //re-process the right-most portion of the range
						start: m.end,
						end:   value_range.end,
					}
				}
				var outputs [][]Range
				if len(ranges) > 0 || (Range{}) != re_process {
					fmt.Printf("Category %s range ", category)
					fmt.Print(value_range)
					fmt.Printf(" became %s ranges: ", k)
					fmt.Print(ranges)
					if (Range{}) != re_process {
						fmt.Print(" and ")
						fmt.Print(re_process)
						fmt.Print(" still needs processing")
					}
					fmt.Println()
				}
				if (Range{}) != re_process {
					outputs = append(outputs, converter(category, re_process))
				}
				for _, r := range ranges {
					outputs = append(outputs, converter(k, r))
				}
				if len(outputs) > 0 {
					return concatMultipleSlices[Range](outputs)
				}
			}
			return converter(k, value_range)
		}
		panic("hit end of recursive function somehow")
	}

	for _, m := range mappings {
		for _, mapping := range m {
			sort.Slice(mapping, func(i, j int) bool {
				return mapping[i].start < mapping[j].start
			})
		}
	}

	var locations []Range
	fmt.Println()
	for _, s := range initialSeeds {
		locations = append(locations, converter(category, s)...)
		fmt.Println("-----------------------")
	}
	fmt.Println(locations)
	sort.Slice(locations, func(i, j int) bool {
		return locations[i].start < locations[j].start
	})
	fmt.Println(locations)
	fmt.Println(locations[0])
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
