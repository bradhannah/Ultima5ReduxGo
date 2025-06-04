package references

import (
	"fmt"
)

/* =========================================================================
   TalkScript AST – updated opcode names per latest spec
   -------------------------------------------------------------------------
   The user requested specific labels for the most common commands. Those
   names have been adopted below while the remaining op‑codes from Ultima V
   are preserved under their original identifiers.  Update or extend as you
   encounter additional bytes in the TLK data.
   =========================================================================*/

type TalkCommand byte

const (
	/* --- Basic printable & prompts ------------------------------------ */
	PlainString TalkCommand = 0x00

	UserInputNotRecognized         TalkCommand = 0x7E
	PromptUserForInputUserInterest TalkCommand = 0x7F
	PromptUserForInputNpcQuestion  TalkCommand = 0x80

	/* --- Substitutions & flow (renamed per request) ------------------- */
	AvatarsName     TalkCommand = 0x81
	EndConversation TalkCommand = 0x82
	Pause           TalkCommand = 0x83
	JoinParty       TalkCommand = 0x84
	GoldPrompt      TalkCommand = 0x85 // was Gold
	Change          TalkCommand = 0x86
	OrBranch        TalkCommand = 0x87 // was Or
	AskName         TalkCommand = 0x88
	KarmaPlusOne    TalkCommand = 0x89
	KarmaMinusOne   TalkCommand = 0x8A
	CallGuards      TalkCommand = 0x8B
	IfElseKnowsName TalkCommand = 0x8C
	NewLine         TalkCommand = 0x8D
	Rune            TalkCommand = 0x8E
	KeyWait         TalkCommand = 0x8F

	StartLabelDef TalkCommand = 0x90 // renamed, was StartLabelDefinition
	// label bytes 0x91‑0x9B represent data (labels 0‑9)

	Label1  TalkCommand = 0x91
	Label2  TalkCommand = 0x92
	Label3  TalkCommand = 0x93
	Label4  TalkCommand = 0x94
	Label5  TalkCommand = 0x95
	Label6  TalkCommand = 0x96
	Label7  TalkCommand = 0x97
	Label8  TalkCommand = 0x98
	Label9  TalkCommand = 0x99
	Label10 TalkCommand = 0x9A

	EndScript       TalkCommand = 0x9F
	StartNewSection TalkCommand = 0xA2

	/* --- Engine‑internal payload codes -------------------------------- */
	ExtortionAmount      TalkCommand = 0xE0
	GoToJail             TalkCommand = 0xE1
	PayGenericExtortion  TalkCommand = 0xE2
	PayHalfGoldExtortion TalkCommand = 0xE3
	MakeAHorse           TalkCommand = 0xE4

	/* --- Branch / label ops ------------------------------------------ */
	GotoLabel                       TalkCommand = 0xFD
	DefineLabel                     TalkCommand = 0xFE
	DoNothingSection                TalkCommand = 0xFF
	PromptUserForInput_NPCQuestion  TalkCommand = 0x80
	PromptUserForInput_UserInterest TalkCommand = 0x7F
)

// String returns a mnemonic for debugging.
func (tc TalkCommand) String() string {
	switch tc {
	case PlainString:
		return "PlainString"
	case UserInputNotRecognized:
		return "UserInputNotRecognized"
	case PromptUserForInputUserInterest:
		return "PromptUserForInputUserInterest"
	case PromptUserForInputNpcQuestion:
		return "PromptUserForInputNpcQuestion"
	case AvatarsName:
		return "AvatarsName"
	case EndConversation:
		return "EndConversation"
	case Pause:
		return "Pause"
	case JoinParty:
		return "JoinParty"
	case GoldPrompt:
		return "GoldPrompt"
	case Change:
		return "Change"
	case OrBranch:
		return "OrBranch"
	case AskName:
		return "AskName"
	case KarmaPlusOne:
		return "KarmaPlusOne"
	case KarmaMinusOne:
		return "KarmaMinusOne"
	case CallGuards:
		return "CallGuards"
	case IfElseKnowsName:
		return "IfElseKnowsName"
	case NewLine:
		return "NewLine"
	case Rune:
		return "Rune"
	case KeyWait:
		return "KeyWait"
	case StartLabelDef:
		return "StartLabelDef"
	case EndScript:
		return "EndScript"
	case StartNewSection:
		return "StartNewSection"
	case ExtortionAmount:
		return "ExtortionAmount"
	case GoToJail:
		return "GoToJail"
	case PayGenericExtortion:
		return "PayGenericExtortion"
	case PayHalfGoldExtortion:
		return "PayHalfGoldExtortion"
	case MakeAHorse:
		return "MakeAHorse"
	case GotoLabel:
		return "GotoLabel"
	case DefineLabel:
		return "DefineLabel"
	case DoNothingSection:
		return "DoNothingSection"
	default:
		return fmt.Sprintf("TalkCommand(0x%02X)", byte(tc))
	}
}

/* --------------------------- AST nodes ----------------------------------- */

type ScriptItem struct {
	Cmd                TalkCommand
	Str                string // valid when Cmd == PlainString
	Num                int    // generic numeric payload
	ItemAdditionalData int    // specialised payload for Change / Gold / etc.
}

type ScriptLine []ScriptItem

type scriptQuestionAnswer struct {
	Questions []string
	Answer    ScriptLine
}

type scriptTalkLabel struct {
	Num            int
	Initial        ScriptLine
	DefaultAnswers []ScriptLine
	QA             []scriptQuestionAnswer
}

type TalkScript struct {
	Lines     []ScriptLine
	Questions map[string]*scriptQuestionAnswer
	Labels    map[int]*scriptTalkLabel
}

// Ask – simple lowercase lookup.
//func (ts *TalkScript) Ask(q string) (ScriptLine, bool) {
//	if ts.Questions == nil {
//		return nil, false
//	}
//	qa, ok := ts.Questions[strings.ToLower(q)]
//	if !ok {
//		return nil, false
//	}
//	return qa.Answer, true
//}

/* ----------------- Convenience constants (fixed line indices) ------------ */

const (
	TalkScriptConstantsName        = 0
	TalkScriptConstantsDescription = 1
	TalkScriptConstantsGreeting    = 2
	TalkScriptConstantsJob         = 3
	TalkScriptConstantsBye         = 4
)

func (sl ScriptLine) isEndOfLabelSection() bool {
	return len(sl) >= 2 &&
		sl[0].Cmd == StartLabelDef &&
		sl[1].Cmd == EndScript
}

func (sl ScriptLine) isLabelDefinition() bool {
	return len(sl) >= 2 &&
		sl[0].Cmd == StartLabelDef &&
		sl[1].Cmd == DefineLabel
}
