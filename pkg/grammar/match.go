package grammar

import "strings"

type Match interface {
	PartiallyMatches(string) bool
	GetSuffixHint(currentStr string) string
}

type StringMatch struct {
	str string
}

type IntMatch struct {
	intMin int
	intMax int
}

func NewStringMatch(str string) StringMatch {
	match := StringMatch{}
	match.str = str
	return match
}

func (m StringMatch) PartiallyMatches(str string) bool {
	return strings.HasPrefix(m.str, str)
}
func (m StringMatch) GetSuffixHint(currentStr string) string {
	if !m.PartiallyMatches(currentStr) {
		// don't give a hint because it doesn't match in the first place
		return ""
	}

	// return the second half after the matched prefix
	return strings.TrimPrefix(m.str, currentStr)
}

func (m IntMatch) PartiallyMatches(str string) bool {
	return true
}

func (m IntMatch) GetSuffixHint(currentStr string) string {
	return ""
}
