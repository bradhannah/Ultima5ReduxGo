package game_state

import (
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

func (p *PlayerCharacter) GetNameAsString() string {
	return strings.TrimRight(string(p.Name[:]), string(rune(0)))
}

func (p *PlayerCharacter) GetKeySpriteIndex() indexes.SpriteIndex {
	if p.Class == 'A' {
		return indexes.Avatar_KeyIndex
	}
	if p.Class == 'M' {
		return indexes.Wizard_KeyIndex
	}
	if p.Class == 'F' {
		return indexes.Fighter_KeyIndex
	}
	if p.Class == 'B' {
		return indexes.Bard_KeyIndex
	}
	return 0
}
