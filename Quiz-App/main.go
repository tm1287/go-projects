package main

import (
	"flag"
	"os"
	"fmt"
	"encoding/csv"
	"strings"
	"time"
	"math/rand"
)

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	perm := rand.Perm(len(lines))
	for i, line := range lines {
		ret[perm[i]] = problem {
			q: strings.TrimSpace(line[0]),
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'problem,solution'")
	timeLimit := flag.Int("limit", 30, "time limit for quiz")
	flag.Parse()
	
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Could not open file: %s", *csvFilename)
		os.Exit(1)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Could not parse csv file")
		os.Exit(1)
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: What is %s? ", i+1, p.q)

		answerC := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerC <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored a %d/%d\n", correct, len(problems))
			return
		case answer := <- answerC:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored a %d/%d\n", correct, len(problems))
}