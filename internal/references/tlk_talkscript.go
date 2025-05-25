// parser.go  (same package "references")
package references

import (
	"strings"
)

// ParseNPCBlob converts the raw TLK byte slice for a single NPC into
// a TalkScript that currently contains only plain strings.
func ParseNPCBlob(blob []byte, dict *WordDict) (*TalkScript, error) {
	const eol = 0x00

	var (
		buf      strings.Builder
		lines    []ScriptLine
		currLine ScriptLine
		addPlain = func() {
			if buf.Len() == 0 {
				return
			}
			currLine = append(currLine, ScriptItem{
				Cmd: PlainString,
				Str: buf.String(),
			})
			buf.Reset()
		}
		flushLine = func() {
			addPlain() // make sure trailing text is included
			if len(currLine) > 0 {
				lines = append(lines, currLine)
				currLine = nil // start a fresh line next time
			}
		}
	)

	for _, b := range blob {
		switch {
		case b == eol:
			flushLine()

			// 1) ASCII + 0x80 letters -----------------------------
		case (b >= 0xA5 && b <= 0xDA) ||
			(b >= 0xE1 && b <= 0xFA) ||
			(b >= 0xA0 && b <= 0xA1):
			buf.WriteByte(b - 0x80)

			// 2) **real opcodes** --------------------------------
		case b == byte(AvatarsName) ||
			b == byte(EndConversation) ||
			b == byte(Pause) /* …add any others you’ve defined… */ :
			addPlain()
			currLine = append(currLine, ScriptItem{Cmd: TalkCommand(b)})

			// 3) compressed‑word bytes ----------------------------
		case dict.IsWordByte(b):
			word, _ := dict.Word(b)
			if buf.Len() > 0 && buf.String()[buf.Len()-1] != ' ' {
				buf.WriteByte(' ')
			}
			buf.WriteString(word)
			buf.WriteByte(' ')

			// 4) unknown fall‑through -----------------------------
		default:
			addPlain()
			currLine = append(currLine, ScriptItem{Cmd: TalkCommand(b)})
		}

	}
	flushLine() // final line, if any

	return &TalkScript{
		Lines:     lines,
		Questions: map[string]*ScriptQuestionAnswer{},
		Labels:    map[int]*ScriptTalkLabel{},
	}, nil
}
