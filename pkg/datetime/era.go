package datetime

type Era int

const (
	EarlyEra = iota
	MiddleEra
	LateEra
)

const (
	beginningOfEra1 = 0
	beginningOfEra2 = 10000
	beginningOfEra3 = 30000
)

func GetEraByTurn(nTurn int) Era {
	if nTurn >= beginningOfEra3 {
		return LateEra
	}
	if nTurn >= beginningOfEra2 {
		return MiddleEra
	}
	return EarlyEra
}
