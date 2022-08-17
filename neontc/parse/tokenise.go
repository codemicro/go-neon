package parse

import "fmt"

type tokenSet struct {
	tokens    [][]byte
	numTokens int

	cursor int
}

func (t *tokenSet) Next() []byte {
	if t.cursor+1 == t.numTokens {
		return nil
	}
	t.cursor += 1
	return t.tokens[t.cursor-1]
}

func (t *tokenSet) Rewind() {
	t.cursor -= 1
}

func (t *tokenSet) Peek() []byte {
	if t.cursor+1 == t.numTokens {
		return nil
	}
	return t.tokens[t.cursor+1]
}

func tokens(input []byte) (*tokenSet, error) {
	const (
		seekStart = iota
		seekEnd
	)

	var (
		output [][]byte

		stage = seekStart
		// the first character of the beginning of this segment
		previousPoint int
	)

	for i, b := range input {

		// Opening sequence: {{
		if b == '{' && getAtIndex(input, i+1) == '{' && getAtIndex(input, i-1) != '\\' {
			if stage == seekEnd {
				return nil, fmt.Errorf("new set of opening curlies at position %d when searching for closing curlies", i)
			}

			output = append(output, input[previousPoint:i])
			stage = seekEnd

			previousPoint = i
		}

		// Closing sequence: }}
		if b == '}' && getAtIndex(input, i-1) == '}' && getAtIndex(input, i-2) != '\\' {
			if stage == seekStart {
				return nil, fmt.Errorf("new set of closing curlies at position %d when searching for opening curlies", i)
			}
			output = append(output, input[previousPoint:i+1])
			stage = seekStart

			previousPoint = i + 1
		}
	}

	if stage == seekEnd {
		return nil, fmt.Errorf("unclosed opening curlies starting at position %d", previousPoint)
	}

	// If the input starts or ends with curlies, we will end up with something
	// akin to []string{"", "blah"....} or []string{..."blah", ""} being
	// outputted.

	// The checks below should prevent that

	if addition := input[previousPoint:]; len(addition) != 0 {
		output = append(output, addition)
	}

	if len(output) != 0 && len(output[0]) == 0 {
		output = output[1:]
	}

	return &tokenSet{
		tokens:    output,
		numTokens: len(output),
	}, nil
}
