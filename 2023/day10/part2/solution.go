package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

type Node struct {
	distance  int
	neighbors []*Node
	point     Point
}

type space int

const (
	unknown space = 0
	A       space = -1
	B       space = 1
	inside  space = 10
	outside space = -10
)

type Spot struct {
	r        rune
	enclosed space
	node     *Node
}

type Point struct {
	y int
	x int
}

type direction int

const (
	north     direction = 0
	east      direction = 90
	south     direction = 180
	west      direction = 270
	northeast direction = (0 + 45) % 360
	northwest direction = (360 - 45) % 360
	southeast direction = (south + 45) % 360
	southwest direction = (south + -45) % 360
)

func turn(d direction, deg int) direction {
	return direction((int(d) + deg) % 360)
}
func opposite(d direction) direction {
	return turn(d, 180)
}

func followPipe(r rune, d direction) direction {
	switch string(r) {
	case "|": //| is a vertical pipe connecting north and south.
		if d == north {
			return south
		} else {
			return north
		}
	case "-": //- is a horizontal pipe connecting east and west.
		if d == east {
			return west
		} else {
			return east
		}
	case "L": //L is a 90-degree bend connecting north and east.
		if d == north {
			return east
		} else {
			return north
		}
	case "J": //J is a 90-degree bend connecting north and west.
		if d == north {
			return west
		} else {
			return north
		}
	case "7": //7 is a 90-degree bend connecting south and west.
		if d == south {
			return west
		} else {
			return south
		}
	case "F": //F is a 90-degree bend connecting south and east.
		if d == south {
			return east
		} else {
			return south
		}
	case ".": //. is ground; there is no pipe in this tile.
		log.Fatalf("loop is not closed, pipe connected to ground\n")
	case "S": //S is the starting Point of the animal
		fmt.Println("Hit start again")
	default:
		log.Fatalln("hit default case unexpectedly")
	}
	panic("oops")
}

func transform(d direction, p Point) Point {
	switch d {
	case north:
		return Point{x: p.x, y: p.y - 1}
	case east:
		return Point{x: p.x + 1, y: p.y}
	case south:
		return Point{x: p.x, y: p.y + 1}
	case west:
		return Point{x: p.x - 1, y: p.y}
	case northeast:
		return Point{x: p.x + 1, y: p.y - 1}
	case northwest:
		return Point{x: p.x - 1, y: p.y - 1}
	case southeast:
		return Point{x: p.x + 1, y: p.y + 1}
	case southwest:
		return Point{x: p.x - 1, y: p.y + 1}
	default:
		return p
	}
}

func allowedDirections(r rune) []direction {
	switch r {
	case '|': //| is a vertical pipe connecting north and south.
		return []direction{north, south}
	case '-': //- is a horizontal pipe connecting east and west.
		return []direction{east, west}
	case 'L': //L is a 90-degree bend connecting north and east.
		return []direction{north, east}
	case 'J': //J is a 90-degree bend connecting north and west.
		return []direction{north, west}
	case '7': //7 is a 90-degree bend connecting south and west.
		return []direction{south, west}
	case 'F': //F is a 90-degree bend connecting south and east.
		return []direction{east, south}
	case '.': //. is ground; there is no pipe in this tile.
		log.Fatalf("loop is not closed, pipe connected to ground\n")
	case 'S': //S is the starting Point of the animal
		fmt.Println("Hit start again")
	default:
		log.Fatalln("hit default case unexpectedly")
	}
	panic("oops")
}

func checkLimits(p Point, matrix [][]Spot) bool {
	return p.x >= 0 && p.y >= 0 && p.y < len(matrix) && p.x < len(matrix[p.y])
}

func isPipe(s Spot) bool {
	return s.node != nil && s.node.distance != -1
}

func printSpot(s Spot) {
	if isPipe(s) {
		fmt.Print(string(s.r))
	} else {
		switch s.enclosed {
		case A:
			fmt.Print("A")
		case B:
			fmt.Print("B")
		case inside:
			fmt.Print("I")
		case outside:
			fmt.Print("O")
		case unknown:
			fmt.Print("U")
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
	var nodes [][]Spot
	start := Point{}
	for i := 0; scanner.Scan(); i++ {
		nodes = append(nodes, make([]Spot, len(scanner.Text())))
		for j, c := range scanner.Text() {
			nodes[i][j] = Spot{r: c}
			if string(c) == "S" {
				nodes[i][j].node = &Node{}
				start = Point{
					y: i,
					x: j,
				}
			}
		}
	}
	var loopDirections []direction // track which directions the start position is attached to so we can replace it w/ it's pipe
	var traverse func(p Point, from *Node, d direction) (*Node, bool)
	traverse = func(p Point, from *Node, d direction) (*Node, bool) {
		if nodes[p.y][p.x].node != nil {
			if string(nodes[p.y][p.x].r) == "S" {
				loopDirections = append(loopDirections, d)
				nodes[p.y][p.x].node.neighbors = append(nodes[p.y][p.x].node.neighbors, from)
				return nodes[p.y][p.x].node, true
			}
			return nodes[p.y][p.x].node, false
		} else {
			nodes[p.y][p.x].node = &Node{distance: -1, point: Point{x: p.x, y: p.y}}
			nodes[p.y][p.x].node.neighbors = []*Node{from}
			toVisit := followPipe(nodes[p.y][p.x].r, d)
			n := transform(toVisit, p)
			if checkLimits(n, nodes) {
				next, connected := traverse(transform(toVisit, p), nodes[p.y][p.x].node, opposite(toVisit))
				if connected {
					nodes[p.y][p.x].node.neighbors = append(nodes[p.y][p.x].node.neighbors, next)
					toLeft := transform(turn(d, 270), p) // get point to the left
					if checkLimits(toLeft, nodes) {
						nodes[toLeft.y][toLeft.x].enclosed = A
					}
					toRight := transform(turn(d, 90), p) // get point to the left
					if checkLimits(toRight, nodes) {
						nodes[toRight.y][toRight.x].enclosed = B
					}
				}
				return nodes[p.y][p.x].node, connected
			}
			return nodes[p.y][p.x].node, false
		}
	}
	root := nodes[start.y][start.x].node
	for _, d := range []direction{north, east, south, west} {
		p := transform(d, start)
		if checkLimits(p, nodes) {
			if string(nodes[p.y][p.x].r) != "." && slices.Contains(allowedDirections(nodes[p.y][p.x].r), opposite(d)) {
				next, connected := traverse(p, root, opposite(d))
				if connected {
					loopDirections = append(loopDirections, d)
					root.neighbors = append(root.neighbors, next)
				}
			}
		}
	}
	slices.Reverse(loopDirections) // reverse directions to match the natural N-E-S-W procession order
	for _, p := range []rune{'|', '-', 'L', 'J', '7', 'F'} {
		if slices.Equal(allowedDirections(p), loopDirections) {
			nodes[start.y][start.x].r = p
			break
		}

	}
	if len(root.neighbors) == 2 {
		positions := root.neighbors
		distance := 1
		visited := []*Node{root}
		for positions[0] != positions[1] {
			for i := 0; i < 2; i++ {
				n := positions[i]
				if !slices.Contains(visited, n) {
					visited = append(visited, n)
					n.distance = distance
					if !slices.Contains(visited, n.neighbors[1]) {
						positions[i] = n.neighbors[1]
					} else if !slices.Contains(visited, n.neighbors[0]) {
						positions[i] = n.neighbors[0]
					}
				}
			}
			distance++
		}
		positions[0].distance = distance
		fmt.Printf("part1: %d\n", distance)
		// time for part 2 - taking a scan-line approach
		area := 0
		for y := range nodes {
			crossings := 0
			numUp := 0
			numDown := 0
			for x := range nodes[y] {
				if isPipe(nodes[y][x]) { //if opposite of A/B (aka which is inside)
					switch string(nodes[y][x].r) {
					case "|": //| is a vertical pipe connecting north and south.
						crossings++
					case "L": //L is a 90-degree bend connecting north and east.
						fallthrough
					case "J": //J is a 90-degree bend connecting north and west.
						numUp++
						if numDown == numUp { //if these are equal, we've completed a single crossing
							crossings++
							numDown = 0
							numUp = 0
						} else if numUp%2 == 0 { //if this is even, we haven't crossed anything
							numUp = 0
						}
					case "F": //F is a 90-degree bend connecting south and east.
						fallthrough
					case "7": //7 is a 90-degree bend connecting south and west.
						numDown++
						if numDown == numUp { //if these are equal, we've completed a single crossing
							crossings++
							numDown = 0
							numUp = 0
						} else if numDown%2 == 0 { //if this is even, we haven't crossed anything
							numDown = 0
						}
					}
				} else if crossings%2 == 0 {
					nodes[y][x].enclosed = outside
				} else {
					nodes[y][x].enclosed = inside
					area++
				}
				printSpot(nodes[y][x])
			}
			fmt.Println()
		}
		fmt.Printf("part2: %d\n", area)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
