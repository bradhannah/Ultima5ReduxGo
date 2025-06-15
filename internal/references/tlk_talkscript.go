// parser.go  (same package "references")
package references

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	minLabelByte = 0x91 // first label byte in a .tlk file
	maxLabelByte = 0x9B // last  (exclusive)
	totalLabels  = 0x0A // ⇐ here – exactly 10 labels (0‑9)
)

// parseNPCBlob converts the raw TLK byte slice for a single NPC into
// a TalkScript that currently contains only plain strings.
func parseNPCBlob(blob []byte, dict *WordDict) (*TalkScript, error) {
	const eol = 0x00

	var (
		buf      strings.Builder
		lines    []ScriptLine
		currLine ScriptLine
		addPlain = func() {
			// sometimes there is zero bytes like Stephen from LB Castle Greeting - but we must add it
			// to maintain index positions
			if buf.Len() == 0 && len(lines) != TalkScriptConstantsGreeting {
				return
				//runtime.Breakpoint()
				//return
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
	const upperHalfHighNunmberSubtractor = 0x80
	for i := 0; i < len(blob); i++ {
		b := blob[i]

		switch {
		//StartLabelDefinition
		//StartLabelDef
		case b == byte(DefineLabel):
			//case b == byte(StartLabelDef):
			addPlain()
			// next byte is the label number (0‑9); consume it
			// guard against running past EOF
			if i+1 >= len(blob) {
				return nil, fmt.Errorf("truncated DefineLabel at end of blob")
			}
			labelNum := int(blob[i+1]) - minLabelByte
			labelNum = 0
			i++ // skip the payload byte
			currLine = append(currLine, ScriptItem{
				Cmd: DefineLabel,
				Num: labelNum,
			})
		case b == eol:
			flushLine()

		// 1) ASCII 0x80 letters -----------------------------
		case (b >= 0xA5 && b <= 0xDA) ||
			(b >= 0xE1 && b <= 0xFA) ||
			(b >= 0xA0 && b <= 0xA1):
			buf.WriteByte(b - upperHalfHighNunmberSubtractor)

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

	ts := TalkScript{
		Lines:     lines,
		Questions: nil,
		Labels:    nil,
	}
	err := ts.BuildIndices()

	if err != nil {
		return nil, err
	}

	return &ts, nil
}

// GetScriptLine returns the raw ScriptLine at idx (false if OOB)
func (ts *TalkScript) GetScriptLine(idx int) (ScriptLine, bool) {
	if idx < 0 || idx >= len(ts.Lines) {
		return nil, false
	}
	return ts.Lines[idx], true
}

func (ts *TalkScript) GetScriptLineLabelIndex(labelNum int) int {
	for idx, line := range ts.Lines {
		if line.isLabelDefinition() && line[1].Num == labelNum {
			return idx
		}
	}
	return -1
}

func (ts *TalkScript) ensure(keys []string, lineIdx int) {
	if lineIdx >= len(ts.Lines) {
		return
	}
	sqa := &scriptQuestionAnswer{Questions: keys, Answer: ts.Lines[lineIdx]}
	for _, k := range keys {
		ts.Questions[k] = sqa
	}
	if lineIdx == TalkScriptConstantsName {
		sqa.Answer[0].Str = fmt.Sprintf("My name is %s", sqa.Answer[0].Str)
	}
}

// BuildIndices transforms the raw Lines slice (produced by parseNPCBlob)
// into fast lookup collections for questions and label jumps.
func (ts *TalkScript) BuildIndices() error {
	if ts == nil {
		return errors.New("nil TalkScript")
	}
	// guard against double‑build
	if ts.Questions != nil && ts.Labels != nil {
		return nil
	}

	ts.Questions = map[string]*scriptQuestionAnswer{}
	ts.Labels = map[int]*scriptTalkLabel{}

	// 1) default hard‑wired Q&A lines -----------------------------

	// Bye always ends the conversation so append explicit opcode
	if TalkScriptConstantsBye < len(ts.Lines) {
		bye := &ts.Lines[TalkScriptConstantsBye]
		*bye = append(*bye, ScriptItem{Cmd: EndConversation})
		_ = ""
	}

	ts.ensure([]string{"name"}, TalkScriptConstantsName)
	ts.ensure([]string{"job", "work"}, TalkScriptConstantsJob)
	ts.ensure([]string{"bye"}, TalkScriptConstantsBye)

	/* ------------------------------------------------------------
	   2) dynamic Q&A section until the first StartLabelDef
	   ----------------------------------------------------------*/
	nIndex := TalkScriptConstantsBye + 1
	simpleQuestions, nIndex := ts.buildQuestions(nIndex)
	MergeMaps(ts.Questions, simpleQuestions)
	//ts.Questions, nIndex = ts.buildQuestions(nIndex)

	/* ------------------------------------------------------------
	   3) label section
	   ----------------------------------------------------------*/
	for nIndex < len(ts.Lines) {
		startLine := ts.Lines[nIndex]

		if startLine.isEndOfLabelSection() {
			if nIndex != (len(ts.Lines) - 1) {
				log.Fatalf("Unexpected end of label section")
			}
			break // no labels – done
		}

		if !startLine.isLabelDefinition() {
			return fmt.Errorf("malformed label start at line %d", nIndex)
		}

		/* ----- gather lines until next StartLabelDef / EndScript ----*/
		for nIndex < len(ts.Lines) {
			line := ts.Lines[nIndex]

			if line[1].Cmd == EndScript {
				// we don't actually save the end of the script
				// I think...
				break
			}

			labelNum := int(line[1].Cmd) - minLabelByte

			label := &scriptTalkLabel{Num: labelNum, Initial: line}
			ts.Labels[labelNum] = label

			// it's a single line only, so we skip this tom foolery below
			if (nIndex+1) >= len(ts.Lines) || ts.Lines[nIndex+1][0].Cmd == StartLabelDef {
				// do nothing - the ScriptTalkLabel will simply have no DefaultAnswer indicating that only the primary
				// label line is read

				// The first example of this happening is Treanna (https://wiki.ultimacodex.com/wiki/Ultima_V_transcript#Treanna)
				// in Label #4 there is no default answer
				nIndex++
				continue
			}

			if line.isEndOfLabelSection() {
				fmt.Sprint("oof")
				break
			}

			label.DefaultAnswers = append(label.DefaultAnswers, ts.Lines[nIndex+1])
			// Advance the Index past the Initial Definition (+0)
			// and past DefaultAnswers (+1).
			// Since we are in this section we are confident that there is more data to consume
			// because we would have exitted the loop above otherwise
			// So we expect question and answers sections
			nIndex += 2

			if label.QA == nil {
				label.QA = make(map[string]*scriptQuestionAnswer)
			}
			var qa map[string]*scriptQuestionAnswer
			qa, nIndex = ts.buildQuestions(nIndex)
			MergeMaps(label.QA, qa)

			fmt.Sprint("oof")
		}
	}

	return nil
}

func (ts *TalkScript) buildQuestions(idx int) (map[string]*scriptQuestionAnswer, int) {
	sqaMap := make(map[string]*scriptQuestionAnswer, 0)
	for idx < len(ts.Lines) {
		line := ts.Lines[idx]
		if len(line) == 0 {
			idx++

			continue
		}

		if line[0].Cmd == StartLabelDef {
			break // jump to label processing
		}

		qStrings := []string{toKey(line[0].Str)}
		// gather chained <OrBranch>
		for idx+1 < len(ts.Lines) && ts.Lines[idx+1].Contains(OrBranch) {
			idx += 2
			qStrings = append(qStrings, toKey(ts.Lines[idx][0].Str))
		}

		if idx+1 >= len(ts.Lines) {
			log.Fatalf("question without answer at line %d", idx)
		}

		answer := ts.Lines[idx+1]
		sqa := &scriptQuestionAnswer{Questions: qStrings, Answer: answer}

		for _, k := range qStrings {
			//if _, exists := ts.Questions[k]; !exists {
			sqaMap[k] = sqa
			//}
		}

		idx += 2
	}

	return sqaMap, idx
}

func toKey(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

type SplitScriptLine = ScriptLine

// IsQuestion heuristic – 1‑6 chars, no spaces.
func (si ScriptItem) IsQuestion() bool {
	if si.Cmd != PlainString {
		return false
	}
	trimmed := strings.TrimSpace(si.Str)
	return len(trimmed) >= 1 && len(trimmed) <= 6 && !strings.Contains(trimmed, " ")
}

// Contains returns true if the line includes the given opcode.
func (sl ScriptLine) Contains(cmd TalkCommand) bool {
	for _, it := range sl {
		if it.Cmd == cmd {
			return true
		}
	}
	return false
}

// SplitIntoSections replicates the intricate splitting logic from the C#
// implementation.  It walks the op‑codes in the line and divides them into
// logical blocks separated by <A2>, label definitions, If/Else branches,
// Change/Gold opcode payloads, etc.  The resulting slice always contains at
// least one entry.
//
//nolint:cyclop
func (sl ScriptLine) SplitIntoSections() []SplitScriptLine {
	// early‑out for the common case: a plain string with no structural
	// op‑codes – just return a single section containing the full line.
	simple := true

	for _, it := range sl {
		switch it.Cmd {
		case PlainString:
			// still simple
		case StartNewSection, IfElseKnowsName, DoNothingSection, DefineLabel,
			Change, GoldPrompt, StartLabelDef:
			simple = false
			break
		default:
			// op‑codes that don’t affect sectioning – ignore
		}
	}
	if simple {
		return []SplitScriptLine{sl}
	}

	var (
		sections       []SplitScriptLine
		nSection       = -1
		first          = true
		forceSplitNext = false
		ensureSection  = func() {
			if nSection < 0 || nSection >= len(sections) {
				sections = append(sections, SplitScriptLine{})
			}
		}
		startNew = func() {
			sections = append(sections, SplitScriptLine{})
			nSection++
		}
	)

	// guarantee at least one section so indices are valid
	startNew()

	for i := 0; i < len(sl); i++ {
		item := sl[i]

		switch item.Cmd {
		case StartNewSection:
			// <A2> – begin new section
			startNew()

		case IfElseKnowsName, DoNothingSection, DefineLabel:
			// stand‑alone section containing only the control opcode
			startNew()
			sections[nSection] = append(sections[nSection], item)
			forceSplitNext = true

		case Change:
			// CHANGE is followed by an item id (as an opcode). We keep them
			// together in their own section so that the caller can inspect
			// item.ItemAdditionalData later if desired.
			startNew()
			if i+1 < len(sl) {
				item.ItemAdditionalData = int(sl[i+1].Cmd)
			}
			sections[nSection] = append(sections[nSection], item)
			i++ // skip payload byte
			forceSplitNext = true

		case GoldPrompt:
			// GOLD is followed by a 3‑char number encoded as PlainString.
			startNew()
			if i+1 < len(sl) && len(sl[i+1].Str) >= 3 {
				digits := sl[i+1].Str[:3]
				var amt int
				fmt.Sscanf(digits, "%d", &amt)
				item.ItemAdditionalData = amt
			}
			sections[nSection] = append(sections[nSection], item)
			i++ // skip payload
			forceSplitNext = true

		case StartLabelDef:
			// must be followed by DefineLabel – keep both together
			startNew()
			sections[nSection] = append(sections[nSection], item)
			if i+1 < len(sl) {
				sections[nSection] = append(sections[nSection], sl[i+1])
				i++ // skip DefineLabel
			}
			forceSplitNext = true

		default:
			if first {
				// first real opcode goes in section 0
				nSection = 0
			}
			if forceSplitNext {
				forceSplitNext = false
				startNew()
			}
			ensureSection()
			sections[nSection] = append(sections[nSection], item)
		}

		first = false
	}

	return sections
}

func MergeMaps[K comparable, V any](dst, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
}
