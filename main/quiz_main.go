package main

import (
	"flag"
	"github.com/anthoturc/golang-playground/main/quiz"
)

func main() {

	csvFileName := flag.String("csv_file", "problems.csv", "The csv quiz file.")
	duration := flag.Int("duration", 30, "The duration of the quiz.")
	flag.Parse()

	quiz.RunTimedQuiz(*csvFileName, *duration)

}
