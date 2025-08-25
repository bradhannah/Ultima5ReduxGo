package widgets

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
)

type DialogStack struct {
	Dialogs []*Widget
}

func (d *DialogStack) DoModalInputBox(question string, textCommand *grammar.TextCommand, keyboard *input.Keyboard) {
	inputBox := NewInputBox(question, textCommand, keyboard)
	d.PushModalDialog(inputBox)
}

func (d *DialogStack) PushModalDialog(dialog Widget) {
	d.Dialogs = append(d.Dialogs, &dialog)
}

func (d *DialogStack) PopModalDialog() *Widget {
	w := d.Dialogs[len(d.Dialogs)-1]
	d.Dialogs = d.Dialogs[:len(d.Dialogs)-1]
	return w
}

func (d *DialogStack) PeekTopModalDialog() *Widget {
	if !d.HasOpenDialog() {
		return nil
	}
	return d.Dialogs[len(d.Dialogs)-1]
}

func (d *DialogStack) HasOpenDialog() bool {
	return len(d.Dialogs) > 0
}

func (d *DialogStack) HasWidgetTypeOnTop(iface interface{}) bool {
	if !d.HasOpenDialog() {
		return false
	}
	return helpers.IsOfTypeInterface(*d.Dialogs[len(d.Dialogs)-1], iface)
}

func (d *DialogStack) GetIndexOfDialogType(iface interface{}) int {
	for i, w := range d.Dialogs {
		if helpers.IsOfTypeInterface(w, iface) {
			return i
		}
	}
	return -1
}

func (d *DialogStack) GetOrAssertTopInputBox() *InputBox {
	if !d.HasOpenDialog() || !d.HasWidgetTypeOnTop(InputBox{}) {
		return nil
	}
	w := *d.Dialogs[len(d.Dialogs)-1]
	ib, ok := (w).(*InputBox)
	if !ok {
		log.Fatal("Unexpected - should find input dialog box") // TODO: CONVERT TO SOFT ERROR - UI recovery possible, should not crash game
	}
	return ib
}

func (d *DialogStack) RemoveWidget(iface interface{}) {
	nIndex := d.GetIndexOfDialogType(iface)
	if nIndex == -1 {
		log.Fatal("Unexpected - should find debug dialog index") // TODO: CONVERT TO SOFT ERROR - UI recovery possible, should not crash game
	}
	d.Dialogs = append(d.Dialogs[:nIndex], d.Dialogs[nIndex+1:]...)
}
