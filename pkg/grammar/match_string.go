package grammar

import "strings"

type MatchString struct {
	Str           string
	Description   string
	CaseSensitive bool
}

func (m MatchString) GetString() string {
	return m.Str
}

func (m MatchString) GetPartialMatches(s string) ([]string, error) {
	if s == "" {
		return []string{}, nil
	}
	if strings.HasPrefix(strings.ToUpper(m.Str), strings.ToUpper(s)) {
		return []string{strings.ToUpper(s)}, nil
	}
	return []string{}, nil
}

func (m MatchString) PartiallyMatches(str string) (bool, error) {
	matches, _ := m.GetPartialMatches(str)
	return len(matches) > 0, nil
}

func (m MatchString) GetSuffixHint(currentStr string) string {
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
