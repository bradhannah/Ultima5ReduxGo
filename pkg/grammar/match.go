package grammar

import (
	"fmt"
	"strconv"
	"strings"
)

type Match interface {
	PartiallyMatches(string) (bool, error)
	GetSuffixHint(currentStr string) string
	GetString() string
	//GetAsInt() int
}

type StringMatch struct {
	Str           string
	Description   string
	CaseSensitive bool
}

func (m StringMatch) GetString() string {
	return m.Str
}

//func (m StringMatch) GetAsInt() int {
//	n, err := strconv.Atoi(m.Str)
//	if err != nil {
//		return n
//	}
//	return 0
//}

type IntMatch struct {
	IntMin      int
	IntMax      int
	Description string
}

func (m IntMatch) GetString() string {
	return fmt.Sprintf("%d to %d", m.IntMax, m.IntMax)
}

//func (m IntMatch) GetAsInt() int {
//	n, err := strconv.Atoi(m.)
//	if err != nil {
//		return n
//	}
//	return 0
//}

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
