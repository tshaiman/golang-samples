package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFile := flag.String("csvFile", "problems.csv", "The CSV file in a format of questions,ans format")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	problems := parseLines(lines)
	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.q)
		var answer string
		fmt.Scanf("%s\n", &answer)
		if answer == problem.a {
			correct++
		}
	}
	fmt.Printf("you got %d correct answers out of %d", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}
