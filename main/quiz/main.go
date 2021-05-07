package main

import "flag"

func main() {

	csvFileName := flag.String("csv_file", "problems.csv", "The csv quiz file.")
	duration := flag.Int("duration", 30, "The duration of the quiz.")
	flag.Parse()

	RunTimedQuiz(*csvFileName, *duration)

}
