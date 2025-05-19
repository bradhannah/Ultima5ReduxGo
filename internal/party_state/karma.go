package party_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type Karma byte

func (k Karma) GetDecreasedKarma(dec Karma) Karma {
	return helpers.Max(k-dec, 0)
}
