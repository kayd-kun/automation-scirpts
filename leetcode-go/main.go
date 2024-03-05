package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Question struct {
	Title      string `json:title`
	Difficulty string `json:difficulty`
	QuestionID string `json:questionId`
}

func main() {
	// ! CHANGE URL INTO DIFFERENT TAG
	// ! CHANGE ARRAY INTO DIFFERENT TYPE WHEN GETTING THE DATA
	url := "https://leetcode.com/problems/tag-data/question-tags/array/"
	questionType := "Array"
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close() // caller must close the response body

	data, err := io.ReadAll(resp.Body)
	check(err)

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		panic(err)
	}

	questionsArray := jsonData["questions"]
	var questions []Question
	questionsSlice, ok := questionsArray.([]interface{})
	if !ok {
		panic("questionsArray is not a slice")
	}
	for _, questionData := range questionsSlice {
		var question Question
		questionJson, _ := json.Marshal(questionData)
		json.Unmarshal(questionJson, &question)
		questions = append(questions, question)
		fmt.Println("Question ID:", question.QuestionID)
		fmt.Println("Question Title:", question.Title)
		fmt.Println("Question Difficulty:", question.Difficulty)
		fmt.Println("================")
	}

	var easyQuestions []Question
	for _, question := range questions {
		if string(question.Difficulty) != "Easy" {
			// fmt.Println("Not Easy Question")
			continue
		}

		easyQuestions = append(easyQuestions, question)
	}

	fmt.Printf("Total easy questions of type %s: %d\n", questionType, len(easyQuestions))

	// Randomly divide the slice into two --> one for patterns, one for practise
	splitSlice(easyQuestions, questionType)
}

func splitSlice(allQuestions []Question, questionType string) {

	n := 1
	var patternQuestions []Question
	var tryQuestions []Question

	for _, v := range allQuestions {
		if n%2 == 0 {
			patternQuestions = append(patternQuestions, v)
		} else {
			tryQuestions = append(tryQuestions, v)
		}
		n++
	}

	// Write to file
	patternsFileName := questionType + "-patterns-easy.txt"
	tryFileName := questionType + "-try-easy.txt"

	writeToFile(patternsFileName, patternQuestions)
	writeToFile(tryFileName, tryQuestions)
}

func writeToFile(fileName string, questions []Question) {
	file, err := os.Create(fileName)
	check(err)
	defer file.Close()
	for _, question := range questions {
		line := fmt.Sprintf("%s: %s\n", question.QuestionID, question.Title)
		_, err := file.WriteString(line)
		check(err)
	}
}
