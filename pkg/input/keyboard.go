package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Keyboard struct {
	MillisecondDelayBetweenKeyPresses int
	timeOfLastKeyPress                int64
}

// TryToRegisterKeyPress
// Tries to register a keypress. Returns true if keypress is allowed according to the time passed since
// the previous key press
func (k *Keyboard) TryToRegisterKeyPress() bool {
	nowMilli := time.Now().UnixMilli()
	if nowMilli-k.timeOfLastKeyPress > int64(k.MillisecondDelayBetweenKeyPresses) {
		k.timeOfLastKeyPress = nowMilli
		return true
	}
	return false
}

func (k *Keyboard) IsBoundKeyPressed(boundKeys []ebiten.Key) bool {
	for _, boundKey := range boundKeys {
		if ebiten.IsKeyPressed(boundKey) {
			return true
		}
	}
	return false
}

func (k *Keyboard) SetLastKeyPressedNow() {
	k.timeOfLastKeyPress = time.Now().UnixMilli()
}

func (k *Keyboard) SetAllowKeyPressImmediately() {
	k.timeOfLastKeyPress = 0
}
