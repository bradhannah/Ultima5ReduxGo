package references

import (
	"bytes"
	"errors"
	"strings"
)

// TALK_OFFSET_ADJUST matches Origin’s 0x80 rule for phrase bytes.
const TALK_OFFSET_ADJUST = 0x80

// WordDict holds the compressed‑word table and a fast byte‑lookup map.
type WordDict struct {
	words       []string // index → word
	byteToIndex map[byte]int
}

// NewWordDict builds the lookup from the string list already extracted from
// DATA.OVL (or wherever you store it now).
//
// The original C# code had many “gap” rules. We replicate them once during
// construction and then pay O(1) per lookup later.
func NewWordDict(raw []string) *WordDict {
	d := &WordDict{
		words:       raw,
		byteToIndex: make(map[byte]int, len(raw)),
	}

	off := 0
	off-- // -1
	d.addRange(1, 7, off)

	off-- // -2
	d.addRange(9, 27, off)

	off-- // -3
	d.addRange(29, 49, off)

	off-- // -4
	d.addRange(51, 64, off)

	off-- // -5
	d.addRange(66, 66, off)

	off-- // -6
	d.addRange(68, 69, off)

	off-- // -7
	d.addRange(71, 71, off)

	off -= 4 // -11  (C# did “i -= 4” here)
	d.addRange(76, 129, off)

	return d
}

func (d *WordDict) addRange(start, stop, offset int) {
	for b := start; b <= stop; b++ {
		d.byteToIndex[byte(b)] = b + offset
	}
}

// IsWordByte reports whether this TLK byte refers to a compressed word.
func (d *WordDict) IsWordByte(b byte) bool {
	_, ok := d.byteToIndex[b]

	return ok
}

// Word returns the expanded word for a TLK code byte.
func (d *WordDict) Word(b byte) (string, error) {
	idx, ok := d.byteToIndex[b]
	if !ok || idx < 0 || idx >= len(d.words) {
		return "", errors.New("no compressed word for code byte")
	}
	return d.words[idx], nil
}

// ---------------------------------------------------------------------------
// Utility helpers (optional)
// ---------------------------------------------------------------------------

// ReplaceMerchantString expands every compressed‑word byte found in a raw shop
// dialog line (leaving “variable” placeholders like %, $, etc. intact).
//
// It mirrors your original ReplaceRawMerchantStringsWithCompressedWords().
func (d *WordDict) ReplaceMerchantString(raw string) (string, error) {
	var out bytes.Buffer
	useWord := false

	for i := 0; i < len(raw); i++ {
		c := raw[i]

		switch {
		case d.IsWordByte(c - TALK_OFFSET_ADJUST):
			word, err := d.Word(c - TALK_OFFSET_ADJUST + 1)
			if err != nil {
				return "", err
			}
			if out.Len() > 0 && !strings.HasSuffix(out.String(), " ") {
				out.WriteByte(' ')
			}
			out.WriteString(word)
			useWord = true

		default:
			if useWord && c != ' ' {
				out.WriteByte(' ')
			}
			out.WriteByte(c)
			useWord = false
		}
	}
	return out.String(), nil
}
