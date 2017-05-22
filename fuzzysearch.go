package main

import (
	"strings"
)

// Fuzz is the penalty for a mismatched character.
// The value should be between 0 and 1.
var Fuzz = 0.5

// FuzzySearch returns the confidence level that the given strings match
func FuzzySearch(input, target string) float64 {
	// No input or target
	if input == "" || target == "" {
		return 0
	}

	// Perfect match
	if input == target {
		return 1
	}

	// Case-insensitive matching
	li := strings.ToLower(input)
	lt := strings.ToLower(target)

	// Cached lengths
	il := len(input)
	tl := len(target)

	// Overall mismatch of strings
	fuzziness := 1.0

	// Denotes the start of search space
	x := 0

	// Amount of confidence in the two characters matching
	score := 0.0

	// Amount of confidence during the process of matching
	total := 0.0

	for i := 0; i < tl; i++ {
		// Substring that denotes the remaining content to search
		slice := li[x:]
		// Case-insensitive character matching
		y := strings.IndexByte(slice, lt[i])

		if y == -1 {
			fuzziness += Fuzz
		} else {
			if x == x+y {
				// First index character match and consecutive characters
				score = 0.7
			} else {
				score = 0.1

				if input[x+y-1] == ' ' {
					// Bonus for acronyms as a result of two character matches
					score += 0.8
				}
			}

			if input[x+y] == target[i] {
				// Bonus for same case
				score += 0.1
			}

			// Update total with the score of the character
			total += score
			// Update x to reflect start of search substring
			x += y + 1
		}
	}

	// Amount of confidence of the two strings matching
	// Reduce penalty of long strings
	confidence := 0.5 * (total/float64(il) + total/float64(tl)) / fuzziness

	if (lt[0] == li[0]) && (confidence < 0.85) {
		confidence += 0.15
	}

	return confidence
}
