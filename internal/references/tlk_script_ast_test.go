package references

import (
	"reflect"
	"testing"
)

func TestQuestionGrouping(t *testing.T) {
	ts := &TalkScript{}

	ts.Lines = []ScriptLine{
		{{Cmd: PlainString, Str: "Name"}},   // 0
		{{Cmd: PlainString, Str: "Job"}},    // 1
		{{Cmd: PlainString, Str: "Bye"}},    // 2
		{{Cmd: PlainString, Str: "Music"}},  // 3
		{{Cmd: PlainString, Str: "Spirit"}}, // 4
		{{Cmd: PlainString, Str: "Answer"}}, // 5
	}

	// Simulate BuildIndices logic
	questionToAnswer := map[string]ScriptLine{
		"musi": ts.Lines[5],
		"spir": ts.Lines[5],
		"name": ts.Lines[0],
		"job":  ts.Lines[1],
		"bye":  ts.Lines[2],
	}
	questionGroupsMap := map[string][]string{}
	answerMap := map[string]ScriptLine{}
	for q, answer := range questionToAnswer {
		hash := hashScriptLine(answer)
		questionGroupsMap[hash] = append(questionGroupsMap[hash], q)
		answerMap[hash] = answer
	}
	var questionGroups []QuestionGroup
	for hash, options := range questionGroupsMap {
		questionGroups = append(questionGroups, QuestionGroup{
			Options: options,
			Script:  answerMap[hash],
		})
	}
	ts.QuestionGroups = questionGroups

	// Validate grouping
	found := false
	for _, group := range ts.QuestionGroups {
		if reflect.DeepEqual(group.Options, []string{"musi", "spir"}) {
			found = true
			if group.Script[0].Str != "Answer" {
				t.Errorf("Expected script to be 'Answer', got '%s'", group.Script[0].Str)
			}
		}
	}
	if !found {
		t.Errorf("Did not find grouped options for musi/spir")
	}

	t.Run("ValidateUpdatedQuestionAnswerFormat", func(t *testing.T) {
		ts := &TalkScript{}

		ts.Lines = []ScriptLine{
			{{Cmd: PlainString, Str: "What is your name?"}}, // 0
			{{Cmd: PlainString, Str: "What is your job?"}},  // 1
			{{Cmd: PlainString, Str: "Goodbye"}},            // 2
			{{Cmd: PlainString, Str: "Answer"}},             // 3
		}

		questionToAnswer := map[string]ScriptLine{
			"name": ts.Lines[3],
			"job":  ts.Lines[3],
			"bye":  ts.Lines[2],
		}
		questionGroupsMap := map[string][]string{}
		answerMap = map[string]ScriptLine{}
		for q, answer := range questionToAnswer {
			hash := hashScriptLine(answer)
			questionGroupsMap[hash] = append(questionGroupsMap[hash], q)
			answerMap[hash] = answer
		}
		var questionGroups []QuestionGroup
		for hash, options := range questionGroupsMap {
			questionGroups = append(questionGroups, QuestionGroup{
				Options: options,
				Script:  answerMap[hash],
			})
		}
		ts.QuestionGroups = questionGroups

		// Validate updated grouping
		found = false
		for _, group := range ts.QuestionGroups {
			if reflect.DeepEqual(group.Options, []string{"name", "job"}) {
				found = true
				if group.Script[0].Str != "Answer" {
					t.Errorf("Expected script to be 'Answer', got '%s'", group.Script[0].Str)
				}
			}
		}
		if !found {
			t.Errorf("Expected to find group with options ['name', 'job']")
		}
	})
}
