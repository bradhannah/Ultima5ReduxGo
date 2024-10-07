package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Keyboard struct {
	MillisecondDelayBetweenKeyPresses int
	timeOfLastKeyPress                int64
	lastKeyPressed                    ebiten.Key
}

// TryToRegisterKeyPress
// Tries to register a keypress. Returns true if keypress is allowed according to the time passed since
// the previous key press
func (k *Keyboard) TryToRegisterKeyPress(key ebiten.Key) bool {
	nowMilli := time.Now().UnixMilli()
	if key != k.lastKeyPressed || nowMilli-k.timeOfLastKeyPress > int64(k.MillisecondDelayBetweenKeyPresses) {
		k.lastKeyPressed = key
		k.timeOfLastKeyPress = nowMilli
		return true
	}
	return false
}

func (k *Keyboard) GetBoundKeyPressed(boundKeys *[]ebiten.Key) *ebiten.Key {
	for _, boundKey := range *boundKeys {
		if ebiten.IsKeyPressed(boundKey) {
			return &boundKey
		}
	}
	return nil

}

func (k *Keyboard) IsBoundKeyPressed(boundKeys *[]ebiten.Key) bool {
	return k.GetBoundKeyPressed(boundKeys) != nil
}

func (k *Keyboard) SetLastKeyPressedNow() {
	k.timeOfLastKeyPress = time.Now().UnixMilli()
}

func (k *Keyboard) SetAllowKeyPressImmediately() {
	k.timeOfLastKeyPress = 0
}
