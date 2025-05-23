package grammar

import "strings"

type MatchStringList struct {
	Strings              []string
	Description          string
	CaseSensitive        bool
	SingleCharacterInput bool
}

func (m MatchStringList) ShouldAutofillWithFirstCharacter() bool {
	return m.SingleCharacterInput
}

func (m MatchStringList) GetDescription() string {
	return m.Description
}

func (m MatchStringList) GetPartialMatches(s string) ([]string, error) {
	matches := make([]string, 0)
	if !m.CaseSensitive {
		s = strings.ToUpper(s)
	}
	for _, str := range m.Strings {
		if strings.HasPrefix(strings.ToUpper(str), s) {
			matches = append(matches, str)
		}
	}
	return matches, nil
}

func (m MatchStringList) PartiallyMatches(s string) (bool, error) {
	matches, _ := m.GetPartialMatches(s)
	return len(matches) > 0, nil
}

func (m MatchStringList) GetSuffixHint(currentStr string) string {
	strs := m.getListOfMatches(currentStr)
	if len(*strs) != 1 {
		return ""
	}
	str := (*strs)[0]
	if m.CaseSensitive {
		return strings.TrimPrefix(str, currentStr)
	}
	return strings.TrimPrefix(strings.ToUpper(str), strings.ToUpper(currentStr))
}

func (m MatchStringList) GetString() string {
	return ""
}

func (m MatchStringList) getListOfMatches(currentStr string) *[]string {
	strs := make([]string, 0, len(m.Strings))

	if !m.CaseSensitive {
		currentStr = strings.ToUpper(currentStr)
	}
	for _, str := range m.Strings {
		if strings.HasPrefix(strings.ToUpper(str), currentStr) {
			strs = append(strs, str)
		}
	}
	return &strs
}
