// script_ast.go
// Package references – dialogue abstract‑syntax tree (AST) types.
//
// These structs are intentionally minimal and have no behaviour other than
// a few helper methods.  All parsing / game‑logic will live in other files so
// these remain a clean data model that can round‑trip with encoding/json.
package references

import "fmt"

// ---------------------------------------------------------------------------
// 1. TalkCommand – every opcode that can appear in a TLK script
// ---------------------------------------------------------------------------

type TalkCommand byte

// NOTE: Only a subset is listed here.  Add the rest as you implement later
// phases of the parser.  Keep them contiguous (or document gaps) so loops
// over ranges are possible without allocation.
const (
	PlainString      TalkCommand = 0x00 // printable text
	AvatarsName      TalkCommand = 0x81
	EndConversation              = 0x82
	Pause                        = 0x83
	JoinParty                    = 0x84
	GoldPrompt                   = 0x85
	Change                       = 0x86
	OrBranch                     = 0x87
	AskName                      = 0x88
	NewLine                      = 0x8D
	StartLabelDef                = 0x90
	EndScript                    = 0x9F
	StartNewSection              = 0xA2
	GotoLabel                    = 0xFD
	DefineLabel                  = 0xFE
	DoNothingSection             = 0xFF
)

// String lets fmt.Printf print a mnemonic instead of a raw number.
func (c TalkCommand) String() string {
	switch c {
	case PlainString:
		return "PlainString"
	case AvatarsName:
		return "AvatarsName"
	case EndConversation:
		return "EndConversation"
	case Pause:
		return "Pause"
	case JoinParty:
		return "JoinParty"
	case GoldPrompt:
		return "Gold"
	case Change:
		return "Change"
	case OrBranch:
		return "Or"
	case AskName:
		return "AskName"
	case NewLine:
		return "NewLine"
	case StartLabelDef:
		return "StartLabelDef"
	case EndScript:
		return "EndScript"
	case StartNewSection:
		return "StartNewSection"
	case GotoLabel:
		return "GotoLabel"
	case DefineLabel:
		return "DefineLabel"
	case DoNothingSection:
		return "DoNothingSection"
	default:
		return fmt.Sprintf("TalkCommand(0x%02X)", byte(c))
	}
}

// ---------------------------------------------------------------------------
// 2. Leaf node – ScriptItem
// ---------------------------------------------------------------------------

// ScriptItem represents a single opcode (and its payload) inside a line.
type ScriptItem struct {
	Cmd TalkCommand
	Str string // for PlainString, %,$, etc.
	Num int    // label #, gold amount, item # etc.
}

// ---------------------------------------------------------------------------
// 3. Mid‑level nodes
// ---------------------------------------------------------------------------

type ScriptLine []ScriptItem

type ScriptQuestionAnswer struct {
	Questions []string   // lower‑case trigger words
	Answer    ScriptLine // NPC response
}

type ScriptTalkLabel struct {
	Num            int          // 0‑9
	Initial        ScriptLine   // always printed when jumping here
	DefaultAnswers []ScriptLine // fallback if no QA matched
	QA             []ScriptQuestionAnswer
}

// ---------------------------------------------------------------------------
// 4. Root node – TalkScript
// ---------------------------------------------------------------------------

// TalkScript is the fully‑parsed dialogue tree for one NPC.
type TalkScript struct {
	Lines     []ScriptLine // all raw lines in original order
	Questions map[string]*ScriptQuestionAnswer
	Labels    map[int]*ScriptTalkLabel
}

// Ask provides an ultra‑simple lookup for Phase‑1 testing.
// Later you’ll enhance it (handle labels, karma, etc.).
func (ts *TalkScript) Ask(q string) (ScriptLine, bool) {
	if ts.Questions == nil {
		return nil, false
	}
	qa, ok := ts.Questions[q]
	if !ok {
		return nil, false
	}
	return qa.Answer, true
}
