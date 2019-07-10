package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	csvFile   = flag.String("csvFile", "problems.csv", "The CSV file in a format of questions,ans format")
	timeLimit = flag.Int("timeLimit", 3, "the time limit for the quiz in seconds")
)

type problem struct {
	q string
	a string
}

func main() {
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
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case answer := <-answerCh:
			if answer == problem.a {
				correct++
			}
		}
	}
	fmt.Printf("you got %d correct answers out of %d\n", correct, len(problems))
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
