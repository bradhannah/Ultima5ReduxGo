package grammar

import (
	"strconv"
	"strings"
)

type Match interface {
	PartiallyMatches(string) (bool, error)
	GetSuffixHint(currentStr string) string
}

type StringMatch struct {
	Name          string
	Str           string
	Description   string
	CaseSensitive bool
}

type IntMatch struct {
	IntMin int
	IntMax int
}

func (m StringMatch) PartiallyMatches(str string) (bool, error) {
	if str == "" {
		return false, nil
	}
	return strings.HasPrefix(strings.ToUpper(m.Str), strings.ToUpper(str)), nil
}
func (m StringMatch) GetSuffixHint(currentStr string) string {
	checkMatch, _ := m.PartiallyMatches(currentStr)
	if !checkMatch {
		// don't give a hint because it doesn't match in the first place
		return ""
	}

	// return the second half after the matched prefix
	if m.CaseSensitive {
		return strings.TrimPrefix(m.Str, currentStr)
	}
	return strings.TrimPrefix(strings.ToUpper(m.Str), strings.ToUpper(currentStr))
}

func (m IntMatch) PartiallyMatches(str string) (bool, error) {
	n, err := strconv.Atoi(str)
	if err != nil {
		return false, err
	}
	return n >= m.IntMin && n <= m.IntMax, nil
}

func (m IntMatch) GetSuffixHint(_ string) string {
	return ""
}
