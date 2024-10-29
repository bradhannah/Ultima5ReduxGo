package grammar

import "strings"

type MatchStringList struct {
	Strings       []string
	Description   string
	CaseSensitive bool
}

func (m MatchStringList) GetPartialMatches(s string) ([]string, error) {
	var matches []string = make([]string, 0)
	if !m.CaseSensitive {
		s = strings.ToUpper(s)
	}
	for _, str := range m.Strings {
		if strings.HasPrefix(strings.ToUpper(str), s) {
			matches = append(matches, str)
			//matches = append(matches, m.GetSuffixHint(str))
			//return true, nil
		}
	}
	return matches, nil
}

func (m MatchStringList) PartiallyMatches(s string) (bool, error) {
	matches, _ := m.GetPartialMatches(s)
	return len(matches) > 0, nil
	//if !m.CaseSensitive {
	//	s = strings.ToUpper(s)
	//}
	//for _, str := range m.Strings {
	//	if strings.HasPrefix(strings.ToUpper(str), s) {
	//		return true, nil
	//	}
	//}
	//return false, nil
}

func (m MatchStringList) GetSuffixHint(currentStr string) string {
	// todo: I could return the common parts of the strings that remain
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

	//var matchStr string
	//var nTotal int = 0
	if !m.CaseSensitive {
		currentStr = strings.ToUpper(currentStr)
	}
	for _, str := range m.Strings {
		if strings.HasPrefix(strings.ToUpper(str), currentStr) {
			strs = append(strs, str)
			//matchStr = str
			//nTotal++
		}
	}
	return &strs
}
