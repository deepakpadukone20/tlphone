// Package tlphone (Tulu Phone) is a phonetic algorithm for indexing
// unicode Tulu words by their pronounciation, like Metaphone for English.
// The algorithm generates three Romanized phonetic keys (hashes) of varying
// phonetic proximity for a given Tulu word.
//
// The algorithm takes into account the context sensitivity of sounds, syntactic
// and phonetic gemination, compounding, modifiers, and other known exceptions
// to produce Romanized phonetic hashes of increasing phonetic affinity that are
// faithful to the pronunciation of the original Tulu word.
//
// `key0` = a broad phonetic hash comparable to a Metaphone key that doesn't account
// for hard sounds or phonetic modifiers
//
// `key1` = is a slightly more inclusive hash that accounts for hard sounds
//
// `key2` = highly inclusive and narrow hash that accounts for hard sounds
// and phonetic modifiers
//
// tlphone was created to aid spelling tolerant Tulu word search, but may
// be useful in tasks like spell checking, word suggestion etc.
//
// This is based on MLphone (https://github.com/knadh/knphone/) for Kannada by Kailash Nadh (c) 2019. https://nadh.in
//
// Deepak Padukone Manjunath (c) 2025. https://deepakpadukone.com | License: GPLv3
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
	"ಬ್ಬ": "B",
	"ಪ್ಪ": "P2", "ಮ್ಮ": "M2",
	"ಯ್ಯ": "Y", "ಲ್ಲ": "L2", "ವ್ವ": "V", "ಶ್ಶ": "S1", "ಸ್ಸ": "S",
	"ಳ್ಳ": "L12",
	"ಕ್ಷ": "KS1",
}

var modifiers = map[string]string{
	"ಾ": "", "ಃ": "", "್": "", "ೃ": "R",
	"ಂ": "3", "ಿ": "4", "ೀ": "4", "ು": "5", "ೂ": "5", "ೆ": "6",
	"ೇ": "6", "ೈ": "7", "ೊ": "8", "ೋ": "8", "ೌ": "9", "ൗ": "9",
}

var (
	regexKey0, _     = regexp.Compile(`[1,2,4-9]`)
	regexKey1, _     = regexp.Compile(`[2,4-9]`)
	regexNonTulu, _  = regexp.Compile(`[\P{Tulu}]`)
	regexAlphaNum, _ = regexp.Compile(`[^0-9A-Z]`)
)

// tlphone is the Tulu-phone tokenizer.
type tlphone struct {
	modCompounds  *regexp.Regexp
	modConsonants *regexp.Regexp
	modVowels     *regexp.Regexp
}

// New returns a new instance of the tlPhone tokenizer.
func New() *tlphone {
	var (
		glyphs []string
		mods   []string
		kn     = &tlphone{}
	)

	// modifiers.
	for k := range modifiers {
		mods = append(mods, k)
	}

	// compounds.
	for k := range compounds {
		glyphs = append(glyphs, k)
	}
	kn.modCompounds, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	// consonants.
	glyphs = []string{}
	for k := range consonants {
		glyphs = append(glyphs, k)
	}
	kn.modConsonants, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	// vowels.
	glyphs = []string{}
	for k := range vowels {
		glyphs = append(glyphs, k)
	}
	kn.modVowels, _ = regexp.Compile(`((` + strings.Join(glyphs, "|") + `)(` + strings.Join(mods, "|") + `))`)

	return kn
}

// Encode encodes a unicode Kannada string to its Roman tlPhone hash.
// Ideally, words should be encoded one at a time, and not as phrases
// or sentences.
func (k *tlphone) Encode(input string) (string, string, string) {
	// key2 accounts for hard and modified sounds.
	key2 := k.process(input)

	// key1 loses numeric modifiers that denote phonetic modifiers.
	key1 := regexKey1.ReplaceAllString(key2, "")

	// key0 loses numeric modifiers that denote hard sounds, doubled sounds,
	// and phonetic modifiers.
	key0 := regexKey0.ReplaceAllString(key2, "")

	return key0, key1, key2
}

func (k *tlphone) process(input string) string {
	// Remove all non-tulu characters.
	input = regexNonTulu.ReplaceAllString(strings.Trim(input, ""), "")

	// All character replacements are grouped between { and } to maintain
	// separatability till the final step.

	// Replace and group modified compounds.
	input = k.replaceModifiedGlyphs(input, compounds, k.modCompounds)

	// Replace and group unmodified compounds.
	for k, v := range compounds {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace and group modified consonants and vowels.
	input = k.replaceModifiedGlyphs(input, consonants, k.modConsonants)
	input = k.replaceModifiedGlyphs(input, vowels, k.modVowels)

	// Replace and group unmodified consonants.
	for k, v := range consonants {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace and group unmodified vowels.
	for k, v := range vowels {
		input = strings.ReplaceAll(input, k, `{`+v+`}`)
	}

	// Replace all modifiers.
	for k, v := range modifiers {
		input = strings.ReplaceAll(input, k, v)
	}

	// Remove non alpha numeric characters (losing the bracket grouping).
	return regexAlphaNum.ReplaceAllString(input, "")
}

func (k *tlphone) replaceModifiedGlyphs(input string, glyphs map[string]string, r *regexp.Regexp) string {
	for _, matches := range r.FindAllStringSubmatch(input, -1) {
		for _, m := range matches {
			if rep, ok := glyphs[m]; ok {
				input = strings.ReplaceAll(input, m, rep)
			}
		}
	}
	return input
}
