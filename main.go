package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "Problems csv file in the format of question,answer")
	timeLimit := flag.Int("limit", -1, "Time Limit per question")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if(err != nil) {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}
	lines := readFile(file)
	problems := parseLines(lines)

	if *timeLimit == -1 {
		playQuiz(problems)
	} else {
		playTimedQuiz(problems, *timeLimit)
	}
	
}

func readFile(f *os.File) [][]string {
	file := csv.NewReader(f)
	lines, err := file.ReadAll()

	if err != nil {
		exit("Failed to read given CSV file")
	}
	return lines
}

type problem struct {
	q string
	a string
}

func parseLines(l [][]string ) []problem {
	probs :=  make([]problem, len(l))
	for i, line := range l {
		probs[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return probs
}

func playQuiz(problems []problem) {
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i, p.q)
		var ans string
		fmt.Scanf("%s", &ans)
		if ans == p.a {
			correct++
			fmt.Println("Correct!")
		} else {
			fmt.Println("Incorrect")
		}
	}
	fmt.Printf("You Scored %d out of %d\n", correct, len(problems))
}

func playTimedQuiz(problems []problem, limit int) {
	correct := 0
	timer := time.NewTimer(time.Duration(limit) * time.Second) //Create a timer
	
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i, p.q)
		answerCh := make(chan string) // answer channel
		// Create a go routine and wait for the answer in the ans channel
		go func ()  { 
			var ans string
			fmt.Scanf("%s", &ans)
			answerCh <- ans
		}()

		select { // to handle multiple communicaton operation
		case answer := <-answerCh: // if we receive answer from channel
				if answer == p.a {
					correct++
					fmt.Println("Correct!")
				} else {
					fmt.Println("Incorrect")
				}
		case <-timer.C: // if time is up
			fmt.Printf("You Scored %d out of %d\n", correct, len(problems))
			return
		}
	}
}


func exit(msg string) {
	fmt.Print(msg)
	os.Exit(1)
}


