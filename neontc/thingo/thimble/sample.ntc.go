package thimble

import (
	bbb "fmt"
	ntcKDHhglErfA "html"
	ntccxAljLHjCx "strconv"
	ntcQsdJutDkuO "strings"
)

var Pointless = 42

type PageInputs struct {
	Name []byte
}

func MyPage(inputs *PageInputs) string {
	ntcRijqFUuoGg := new(ntcQsdJutDkuO.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	if true == false {
		_, _ = ntcRijqFUuoGg.WriteString("\n    ok\n    ")
	} else if true != false {
		_, _ = ntcRijqFUuoGg.WriteString("\n    what\n    ")
	} else {
		_, _ = ntcRijqFUuoGg.WriteString("\n    hmmm\n    ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	for _, char := range inputs.Name {
		_, _ = ntcRijqFUuoGg.WriteString("\n        <span>")
		_, _ = ntcRijqFUuoGg.WriteString(ntccxAljLHjCx.FormatUint(uint64(char), 10))
		_, _ = ntcRijqFUuoGg.WriteString("</span>\n    ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    <h1>Hello ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcKDHhglErfA.EscapeString(string(inputs.Name)))
	_, _ = ntcRijqFUuoGg.WriteString("!</h1>\n    <h2>This is ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcKDHhglErfA.EscapeString(bbb.Sprintf("for%s", "matted")))
	_, _ = ntcRijqFUuoGg.WriteString("</h2>\n    <p>This was generated using a Neon template.</p>\n    <p>This is a double bracket: \\{[ \\{{</p>\n    <p>And now in the other direction!: \\]} \\}}</p>\n\n    ")
	print("hi")
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	type MorePageInputs struct {
		Name string
	}

	_, _ = ntcRijqFUuoGg.WriteString("\n\n")
	return ntcRijqFUuoGg.String()
}
