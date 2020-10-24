package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const defaultProblemFileName = "problem.csv"

var problemFileName = "..."

var (
	correctAnswers int
	totalQuestions int
	quizDuration   time.Duration = 10 * time.Second
)

func main() {
	// define flags
	var (
		flagProblemsFilename = flag.String("p", defaultProblemFileName, "the path to the problem CSV file")
		flagTimer            = flag.Duration("t", quizDuration, "The max time for the quiz")
		flagShuffle          = flag.Bool("s", false, "Shuffle the quiz questions")
	)
	flag.Parse()
	if flagProblemsFilename == nil ||
		flagTimer == nil ||
		flagShuffle == nil {
		fmt.Printf("Missing problems filename and/or timer")
		return
	}

	// wait for enter to start the quiz
	fmt.Printf("Hit enter to start the quiz from %q in %v ?",
		*flagProblemsFilename, *flagTimer)
	fmt.Scanln()

	// read the csv file
	f, err := os.Open(*flagProblemsFilename)
	if err != nil {
		fmt.Printf("Failed to open file: %v \n", err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	questions, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read csv file: %v\n", err)
		return
	}
	if *flagShuffle {
		fmt.Println("Shuffling ...")
		rand.Shuffle(len(questions),
			func(i, j int) {
				questions[i], questions[j] = questions[j], questions[i]
			})
	}
	totalQuestions = len(questions)
	quizDone := startQuiz(questions)

	// wait for timer or quiz ends.
	quizTimer := time.NewTimer(*flagTimer).C

	select {
	case <-quizDone:
	case <-quizTimer:
	}

	fmt.Printf("hello after The duration %v, Your result is: %d / %d\n", quizDuration, correctAnswers, totalQuestions)
}

func startQuiz(questions [][]string) chan bool {
	done := make(chan bool)
	go func() {
		for i, record := range questions {
			question, correctAnswer := record[0], record[1]
			fmt.Printf("%d. %s ? \n", i+1, question)
			var answer string
			_, err := fmt.Scan(&answer)
			if err != nil {
				fmt.Printf("failed to scan: %v\n", err)
				return
			}
			// clean up answer
			answer = strings.TrimSpace(answer)
			answer = strings.ToLower(answer)
			if answer == correctAnswer {
				correctAnswers++
			}
		}
		done <- true
	}()
	return done

}
