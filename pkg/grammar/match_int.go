package grammar

import (
	"fmt"
	"strconv"
)

type MatchInt struct {
	IntMin      int
	IntMax      int
	Description string
}

func (m MatchInt) GetDescription() string {
	return m.Description
}

func (m MatchInt) GetPartialMatches(s string) ([]string, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return []string{}, err
	}
	if n >= m.IntMin && n <= m.IntMax {
		return []string{fmt.Sprintf("%d", n)}, nil
	}
	return []string{}, nil
}

func (m MatchInt) GetString() string {
	return fmt.Sprintf("%d to %d", m.IntMax, m.IntMax)
}

func (m MatchInt) PartiallyMatches(str string) (bool, error) {
	n, err := strconv.Atoi(str)
	if err != nil {
		return false, err
	}
	return n >= m.IntMin && n <= m.IntMax, nil
}

func (m MatchInt) GetSuffixHint(_ string) string {
	return ""
}
