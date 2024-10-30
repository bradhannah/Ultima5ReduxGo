package grammar

type Match interface {
	PartiallyMatches(string) (bool, error)
	GetPartialMatches(string) ([]string, error)
	GetSuffixHint(currentStr string) string
	GetString() string
	GetDescription() string
}
