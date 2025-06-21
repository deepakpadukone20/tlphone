// Package tlphone (Tulu Phone) is a phonetic algorithm for indexing
// unicode Tulu words by their pronunciation, like Metaphone for English.
// This adapts the Kannada-based knphone logic for Tulu.
package tlphone

import (
	"regexp"
	"strings"
)

var vowels = map[string]string{
	"ಅ": "A", "ಆ": "A", "ಇ": "I", "ಈ": "I", "ಉ": "U", "ಊ": "U", "ಋ": "R",
	"ಎ": "E", "ಏ": "E", "ಐ": "AI", "ಒ": "O", "ಓ": "O", "ಔ": "O",
}

var consonants = map[string]string{
	"ಕ": "K", "ಖ": "K", "ಗ": "K", "ಘ": "K", "ಙ": "NG",
	"ಚ": "C", "ಛ": "C", "ಜ": "J", "ಝ": "J", "ಞ": "NJ",
	"ಟ": "T", "ಠ": "T", "ಡ": "T", "ಢ": "T", "ಣ": "N1",
	"ತ": "0", "ಥ": "0", "ದ": "0", "ಧ": "0", "ನ": "N",
	"ಪ": "P", "ಫ": "F", "ಬ": "B", "ಭ": "B", "ಮ": "M",
	"ಯ": "Y", "ರ": "R", "ಲ": "L", "ವ": "V",
	"ಶ": "S1", "ಷ": "S1", "ಸ": "S", "ಹ": "H",
	"ಳ": "L1", "ೞ": "Z", "ಱ": "R1",
}

var compounds = map[string]string{
	"ಕ್ಕ": "K2", "ಗ್ಗಾ": "K", "ಙ್ಙ": "NG",
	"ಚ್ಚ": "C2", "ಜ್ಜ": "J", "ಞ್ಞ": "NJ",
	"ಟ್ಟ": "T2", "ಣ್ಣ": "N2",
	"ತ್ತ": "0", "ದ್ದ": "D", "ದ್ಧ": "D", "ನ್ನ": "NN",
	"ಬ್ಬ": "B", "ಪ್ಪ": "P2", "ಮ್ಮ": "M2",
	"ಯ್ಯ": "Y", "ಲ್ಲ": "L2", "ವ್ವ": "V",
	"ಶ್ಶ": "S1", "ಸ್ಸ": "S", "ಳ್ಳ": "L12", "ಕ್ಷ": "KS1",
}

var modifiers = map[string]string{
	"ಾ": "", "ಃ": "", "್": "", "ೃ": "R",
	"ಂ": "3", "ಿ": "4", "ೀ": "4", "ು": "5", "ೂ": "5", "ೆ": "6",
	"ೇ": "6", "ೈ": "7", "ೊ": "8", "ೋ": "8", "ೌ": "9", "ൗ": "9",
}

var (
	regexKey0     = regexp.MustCompile(`[1,2,4-9]`)
	regexKey1     = regexp.MustCompile(`[2,4-9]`)
	regexNonTulu  = regexp.MustCompile(`[\P{Kannada}]`)
	regexAlphaNum = regexp.MustCompile(`[^0-9A-Z]`)
)

type TLPhone struct {
	modCompounds  *regexp.Regexp
	modConsonants *regexp.Regexp
	modVowels     *regexp.Regexp
}

func New() *TLPhone {
	var (
		glyphs []string
		mods   []string
		tl     = &TLPhone{}
	)

	for k := range modifiers {
		mods = append(mods, k)
	}

	for k := range compounds {
		glyphs = append(glyphs, k)
	}
	tl.modCompounds = regexp.MustCompile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	glyphs = []string{}
	for k := range consonants {
		glyphs = append(glyphs, k)
	}
	tl.modConsonants = regexp.MustCompile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	glyphs = []string{}
	for k := range vowels {
		glyphs = append(glyphs, k)
	}
	tl.modVowels = regexp.MustCompile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	return tl
}

func (k *TLPhone) Encode(input string) (string, string, string) {
	key2 := k.process(input)
	key1 := regexKey1.ReplaceAllString(key2, "")
	key0 := regexKey0.ReplaceAllString(key2, "")
	return key0, key1, key2
}

func (k *TLPhone) process(input string) string {
	input = regexNonTulu.ReplaceAllString(strings.TrimSpace(input), "")

	input = k.replaceModifiedGlyphs(input, compounds, k.modCompounds)
	for ck, cv := range compounds {
		input = strings.ReplaceAll(input, ck, `{`+cv+`}`)
	}
	input = k.replaceModifiedGlyphs(input, consonants, k.modConsonants)
	input = k.replaceModifiedGlyphs(input, vowels, k.modVowels)
	for ck, cv := range consonants {
		input = strings.ReplaceAll(input, ck, `{`+cv+`}`)
	}
	for vk, vv := range vowels {
		input = strings.ReplaceAll(input, vk, `{`+vv+`}`)
	}
	for mk, mv := range modifiers {
		input = strings.ReplaceAll(input, mk, mv)
	}

	return regexAlphaNum.ReplaceAllString(input, "")
}

func (k *TLPhone) replaceModifiedGlyphs(input string, glyphs map[string]string, r *regexp.Regexp) string {
	for _, matches := range r.FindAllStringSubmatch(input, -1) {
		for _, m := range matches {
			if rep, ok := glyphs[m]; ok {
				input = strings.ReplaceAll(input, m, rep)
			}
		}
	}
	return input
}
