package thimble

import (
	bbb "fmt"
	ntccxAljLHjCx "html"
	ntcQsdJutDkuO "strings"
)

var Pointless = 42

type PageInputs struct {
	Name string
}

func MyPage(inputs *PageInputs) string {
	ntcRijqFUuoGg := new(ntcQsdJutDkuO.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    <h1>Hello ")
	_, _ = ntcRijqFUuoGg.WriteString(ntccxAljLHjCx.EscapeString(inputs.Name))
	_, _ = ntcRijqFUuoGg.WriteString("!</h1>\n    <h2>This is ")
	_, _ = ntcRijqFUuoGg.WriteString(ntccxAljLHjCx.EscapeString(bbb.Sprintf("for%s", "matted")))
	_, _ = ntcRijqFUuoGg.WriteString("</h2>\n    <p>This was generated using a Neon template.</p>\n    <p>This is a double bracket: \\{[ \\{{</p>\n    <p>And now in the other direction!: \\]} \\}}</p>\n\n    ")
	print("hi")
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	type MorePageInputs struct {
		Name string
	}

	_, _ = ntcRijqFUuoGg.WriteString("\n\n")
	return ntcRijqFUuoGg.String()
}
