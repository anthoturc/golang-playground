package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

// Define a question structure
// that will hold the question (q)
// and answer (a)
type question struct {
	q string
	a string
}

func getQuestionsFromCsv(csvName string) []question {
	f, err := os.Open(csvName)
	if err != nil {
		log.Fatal(err)
	}

	// Read all the data in the csv and
	// marshall each row into a question
	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	questions := make([]question, len(records))
	for i, record := range records {
		// The assumption here is that data
		// in the CSV is in the [question, answer]
		// format
		questions[i] = question{
			q: record[0],
			a: record[1],
		}
	}

	return questions
}

// Quiz without timer
func RunNoTimeLimitQuiz(csvName string) {
	questions := getQuestionsFromCsv(csvName)
	numCorrect := 0

	var actualAnswer string
	for _, question := range questions {

		fmt.Printf("%s = ", question.q)
		fmt.Scanln(&actualAnswer)

		if actualAnswer == question.a {
			numCorrect++
		}
	}

	fmt.Printf("Score: %d/%d\n", numCorrect, len(questions))
}

func RunTimedQuiz(csvName string, duration int) {
	questions := getQuestionsFromCsv(csvName)
	numQuestionsCorrect := 0

	quizTimer := time.NewTimer(time.Duration(duration) * time.Second)

	// The finishedCh channel will enable us to wait on the go routine
	// that is running on the quiz. We need to wait so that we have an
	// accurate score
	finishedCh := make(chan struct{})
	// The answerCh channel will let us know if the user supplied data
	answerCh := make(chan string)

	// Run the quiz in a go routine
	// this will allow us to stop the routine
	// whenever the timer has put data on its
	// channel
	go func(numQuestions int) {
		i := 0
		isQuestionAlreadyAsked := false
		for i < numQuestions {
			select {
			// Check if the timer has stopped
			case <-quizTimer.C:
				finishedCh <- struct{}{}
				return
			// Use a default clause to avoid blocking
			// on the quizTimer channel
			default:
				// To avoid duplicating questions,
				// only ask them once
				if !isQuestionAlreadyAsked {
					fmt.Printf("%s = ", questions[i].q)
					isQuestionAlreadyAsked = true
				}
				select {
				// This 'async' check for input on the user channel
				// will help avoid blocking on user input
				case userInput := <-answerCh:
					if questions[i].a == userInput {
						// This is the only go routine that is modifying
						// this variable so it is safe
						numQuestionsCorrect++
					}
					i++
					isQuestionAlreadyAsked = false
				default:
					// Use a default clause to avoid blocking
					// on the answerCh
				}
			}
		}
		// If the user answers all questions then we still need to signal
		// that we are done because there is a chance that the timer is no
		// longer checked.
		finishedCh <- struct{}{}
	}(len(questions))

	// This routine will let us know where there is user input
	// through the answerCh. There are some paths where we can leak this
	// routine but it shouldn't be a huge issue since its just a single one
	go func() {
		for {
			select {
			case <-finishedCh:
				return
			default:
				var userInput string
				fmt.Scanln(&userInput)
				answerCh <- userInput
			}
		}
	}()

	// We need to show the user the score when
	// a. the timer is done
	// b. the user got through all the questions
	// The go routing responsible for running the test
	// will let us know when we can proceed.
	<-finishedCh
	fmt.Printf("\n\nScore: %d/%d\n", numQuestionsCorrect, len(questions))
}
