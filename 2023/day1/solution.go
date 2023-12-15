package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var numbers_map = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

type Number struct {
	str string
	val string
}

func main() {
	first_chars := map[string][]Number{}
	for k := range numbers_map {
		f := string(k[0])
		first_chars[f] = append(first_chars[f], Number{
			str: k[1:],          // need to drop the first char as it's already been consumed by this step
			val: numbers_map[k], // the number as a single digit char
		})
	}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	sum := 0
	isDigit := regexp.MustCompile(`^[0-9]$`)
	for scanner.Scan() {
		first, last := "", ""
		options := []Number{}
		for _, slice := range strings.Split(scanner.Text(), "") {
			number := ""
			if isDigit.MatchString(slice) {
				number = slice
			} else {
				if len(options) != 0 {
					new_options := []Number{}
					for _, o := range options {
						if slice == string(o.str[0]) {
							if len(o.str) == 1 { // if we matched the last char in the string
								number = o.val
							} else {
								new_options = append(new_options, Number{str: o.str[1:], val: o.val})
							}
						}
					}
					options = new_options
				}
				if possible_words, ok := first_chars[slice]; ok {
					options = append(options, possible_words...)
				}

			}

			if number != "" {
				if first == "" {
					first = number
					last = number
				} else {
					last = number
				}
			}
		}
		num, err := strconv.Atoi(first + last)
		if err == nil {
			fmt.Printf("Adding %d to sum\n", num)
			sum += num
		}
	}
	fmt.Printf("Final Sum is %d\n", sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
