package parse

import "fmt"

type rawToken struct {
	pos  int64
	cont []byte
}

type tokenSet struct {
	tokens    []*rawToken
	numTokens int

	cursor int
}

func (t *tokenSet) Next() *rawToken {
	if t.cursor+1 == t.numTokens {
		return nil
	}
	t.cursor += 1
	return t.tokens[t.cursor-1]
}

func (t *tokenSet) Rewind() {
	t.cursor -= 1
}

func (t *tokenSet) Peek() *rawToken {
	if t.cursor+1 == t.numTokens {
		return nil
	}
	return t.tokens[t.cursor+1]
}

func tokens(fs *FileSet, positionBase int64, input []byte) (*tokenSet, error) {
	const (
		seekStart = iota
		seekEnd
	)

	var (
		output []*rawToken

		stage = seekStart
		// the first character of the beginning of this segment
		previousPoint int
	)

	for i, b := range input {

		// Opening sequence: {{ or {[
		if b == '{' && (getAtIndex(input, i+1) == '{' || getAtIndex(input, i+1) == '[') && getAtIndex(input, i-1) != '\\' {
			if stage == seekEnd {
				return nil, fmt.Errorf("%s: new set of opening brackets when searching for closing brackets", fs.ResolvePosition(positionBase+int64(i)))
			}

			output = append(output, &rawToken{positionBase + int64(previousPoint), input[previousPoint:i]})
			stage = seekEnd

			previousPoint = i
		}

		// Closing sequence: }} or ]}
		if b == '}' && (getAtIndex(input, i-1) == '}' || getAtIndex(input, i-1) == ']') && getAtIndex(input, i-2) != '\\' {
			if stage == seekStart {
				return nil, fmt.Errorf("%s: new set of closing brackets when searching for opening brackets", fs.ResolvePosition(positionBase+int64(i-1)))
			}
			output = append(output, &rawToken{positionBase + int64(previousPoint), input[previousPoint : i+1]})
			stage = seekStart

			previousPoint = i + 1
		}
	}

	if stage == seekEnd {
		return nil, fmt.Errorf("%s: unclosed opening brackets", fs.ResolvePosition(positionBase+int64(previousPoint)))
	}

	// If the input starts or ends with curlies, we will end up with something
	// akin to []string{"", "blah"....} or []string{..."blah", ""} being
	// outputted.

	// The checks below should prevent that

	if addition := input[previousPoint:]; len(addition) != 0 {
		output = append(output, &rawToken{positionBase + int64(previousPoint), addition})
	}

	if len(output) != 0 && len(output[0].cont) == 0 {
		output = output[1:]
	}

	return &tokenSet{
		tokens:    output,
		numTokens: len(output),
	}, nil
}
