package utils

import "testing"

func TestGradeResponseRejectsInvalidOptionIndex(t *testing.T) {
	options := []Option{{Text: "First"}, {Text: "Second"}}
	correctIndex := 0
	answer := &Answer{Text: "First", OptionIndex: &correctIndex}

	for _, index := range []int{-1, -100, len(options), 100} {
		result := GradeResponse("multiple_choice", answer, options, "First", &index, "")
		if result.Error != "Selected option index is out of range" {
			t.Errorf("index %d: expected range error, got %+v", index, result)
		}
		if result.IsCorrect != nil {
			t.Errorf("index %d: invalid response must not be graded", index)
		}
	}

	index := 0
	result := GradeResponse("multiple_choice", answer, nil, "First", &index, "")
	if result.Error == "" || result.IsCorrect != nil {
		t.Fatalf("empty options: expected ungraded validation error, got %+v", result)
	}
}

func TestGradeResponseValidSelections(t *testing.T) {
	options := []Option{{Text: "First"}, {Text: "Second"}}
	correctIndex := 0
	answer := &Answer{Text: "First", OptionIndex: &correctIndex}

	for index, option := range options {
		// The index is authoritative even when supplied text disagrees.
		result := GradeResponse("multiple_choice", answer, options, "Conflicting text", &index, "")
		if result.Error != "" || result.SelectedOption != option.Text || result.IsCorrect == nil {
			t.Fatalf("index %d: unexpected result %+v", index, result)
		}
		if *result.IsCorrect != (index == correctIndex) {
			t.Errorf("index %d: incorrect grading", index)
		}
	}

	result := GradeResponse("multiple_choice", answer, options, "First", nil, "")
	if result.Error != "" || result.IsCorrect == nil || !*result.IsCorrect {
		t.Fatalf("text-only selection should still work: %+v", result)
	}
}

func TestGradeResponseShortAnswer(t *testing.T) {
	result := GradeResponse("short_answer", &Answer{Text: "Answer"}, nil, "", nil, " answer ")
	if result.Error != "" || result.AnswerText != "answer" || result.IsCorrect == nil || !*result.IsCorrect {
		t.Fatalf("short-answer grading changed: %+v", result)
	}
}
