package utils

import (
	"encoding/json"
	"strings"
)

type Option struct {
	Text string `json:"text"`
}

type Answer struct {
	Text        string `json:"text"`
	OptionIndex *int   `json:"optionIndex,omitempty"`
}

func NormalizeOptions(raw []Option) []Option {
	var out []Option
	for _, o := range raw {
		t := strings.TrimSpace(o.Text)
		if t != "" {
			out = append(out, Option{Text: t})
		}
	}
	return out
}

func BuildQuestionFields(qType, prompt string, options []Option, answer *Answer) (options2 []Option, answer2 *Answer, errMsg string) {
	normalized := []Option{}
	if qType == "multiple_choice" {
		normalized = NormalizeOptions(options)
		if len(normalized) < 2 {
			return nil, nil, "Multiple choice questions need at least 2 options"
		}
	}
	if qType == "multiple_choice" && (answer == nil || answer.OptionIndex == nil) {
		return nil, nil, "A valid correct answer (text and optionIndex) is required"
	}
	return normalized, answer, ""
}

func IsMcqCorrect(correctAnswer *Answer, selectedIndex *int, selectedText string) bool {
	if correctAnswer == nil {
		return false
	}
	if selectedIndex != nil && correctAnswer.OptionIndex != nil {
		return *selectedIndex == *correctAnswer.OptionIndex
	}
	return strings.TrimSpace(selectedText) == correctAnswer.Text
}

func MarshalJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

type GradeResult struct {
	SelectedOption      string
	SelectedOptionIndex *int
	AnswerText          string
	IsCorrect           *bool
	Error               string
}

func GradeResponse(qType string, correctAnswer *Answer, options []Option, selectedOption string, selectedOptionIndex *int, answerText string) GradeResult {
	if qType == "multiple_choice" {
		if selectedOptionIndex != nil && *selectedOptionIndex < len(options) {
			selectedOption = options[*selectedOptionIndex].Text
		}
		if selectedOption == "" {
			return GradeResult{Error: "Selected option is required for multiple choice"}
		}
		correct := IsMcqCorrect(correctAnswer, selectedOptionIndex, selectedOption)
		return GradeResult{SelectedOption: selectedOption, SelectedOptionIndex: selectedOptionIndex, IsCorrect: &correct}
	}

	if strings.TrimSpace(answerText) == "" {
		return GradeResult{Error: "Answer text is required for short answer"}
	}
	var correct *bool
	if correctAnswer != nil && correctAnswer.Text != "" {
		match := strings.EqualFold(strings.TrimSpace(answerText), strings.TrimSpace(correctAnswer.Text))
		correct = &match
	}
	return GradeResult{AnswerText: strings.TrimSpace(answerText), IsCorrect: correct}
}
