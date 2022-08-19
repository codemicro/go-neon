package thimble

import (
	"bytes"
)

var Pointless = 42

type PageInputs struct {
	Name string
}

func MyPage(inputs *PageInputs) string {
	ntcRijqFUuoGg := new(bytes.Buffer)
	_, _ = ntcRijqFUuoGg.Write([]byte("\n\n    <h1>Hello "))
	_, _ = ntcRijqFUuoGg.Write([]byte(inputs.Name))
	_, _ = ntcRijqFUuoGg.Write([]byte("!</h1>\n    <p>This was generated using a Neon template.</p>\n    <p>This is a double curly bracket: \\{{</p>\n    <p>And now in the other direction!: \\}}</p>\n\n    "))
	print("hi")
	_, _ = ntcRijqFUuoGg.Write([]byte("\n\n    "))
	type MorePageInputs struct {
		Name string
	}
	_, _ = ntcRijqFUuoGg.Write([]byte("\n\n"))
	return ntcRijqFUuoGg.String()
}
