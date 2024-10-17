package input

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Keyboard struct {
	millisecondDelayBetweenKeyPresses int
	timeOfLastKeyPress                int64
	lastKeyPressed                    ebiten.Key
	forceRegisterNextKeyPressed       bool
}

func NewKeyboard(millisecondDelayBetweenKeyPresses int) *Keyboard {
	k := &Keyboard{}
	k.millisecondDelayBetweenKeyPresses = millisecondDelayBetweenKeyPresses

	return k
}

// TryToRegisterKeyPress
// Tries to register a keypress. Returns true if keypress is allowed according to the time passed since
// the previous key press
func (k *Keyboard) TryToRegisterKeyPress(key ebiten.Key) bool {
	nowMilli := time.Now().UnixMilli()
	if key != k.lastKeyPressed || nowMilli-k.timeOfLastKeyPress > int64(k.millisecondDelayBetweenKeyPresses) || k.forceRegisterNextKeyPressed {
		k.lastKeyPressed = key
		k.timeOfLastKeyPress = nowMilli
		k.forceRegisterNextKeyPressed = false
		return true
	}
	return false
}

func (k *Keyboard) GetMsSinceLastKeyPress() int {
	return int(time.Now().UnixMilli() - k.timeOfLastKeyPress)
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
	//k.timeOfLastKeyPress = 0
	k.forceRegisterNextKeyPressed = true
}

func (k *Keyboard) IsKeyPressedGetKey(key ebiten.Key) *ebiten.Key {
	if ebiten.IsKeyPressed(key) {
		return &key
	}
	return nil
}

func (k *Keyboard) GetKeyLetter(key ebiten.Key) string {
	keyInt := int(key)
	if keyInt >= int(ebiten.KeyA) && keyInt <= int(ebiten.KeyZ) {
		return key.String()
	}
	if keyInt >= int(ebiten.Key0) && keyInt <= int(ebiten.Key9) {
		res := fmt.Sprintf("%d", keyInt-int(ebiten.Key0))
		return res
	}
	if keyInt >= int(ebiten.KeyDigit0) && keyInt <= int(ebiten.KeyDigit9) {
		res := fmt.Sprintf("%d", keyInt-int(ebiten.KeyDigit0))
		return res
	}
	if keyInt >= int(ebiten.KeyNumpad0) && keyInt <= int(ebiten.KeyNumpad9) {
		res := fmt.Sprintf("%d", keyInt-int(ebiten.KeyNumpad0))
		return res
	}
	return key.String()
}

func (k *Keyboard) GetAlphaNumericPressed() (*ebiten.Key, string) {
	var key *ebiten.Key

	for i := int(ebiten.KeyA); i < int(ebiten.KeyZ); i++ {
		key = k.IsKeyPressedGetKey(ebiten.Key(i))
		if key != nil {
			return key, k.GetKeyLetter(*key)
		}
	}
	for i := int(ebiten.Key0); i < int(ebiten.Key9); i++ {
		key = k.IsKeyPressedGetKey(ebiten.Key(i))
		if key != nil {
			return key, k.GetKeyLetter(*key)
		}
	}
	return nil, ""
}
