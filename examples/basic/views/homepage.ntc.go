package views

import (
	ntcbGHNBxxTTK "html"
	"math/rand"
	ntcVqcLceiEsr "strconv"
	ntcYxcNAgvNPP "strings"
)

func Homepage(name string) string {
	ntcRijqFUuoGg := new(ntcYxcNAgvNPP.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n    <html>\n    <head>\n        <title>Hello!</title>\n    </head>\n    <body>\n\n        ")
	var seed int64
	for _, char := range name {
		seed *= 10
		seed += int64(char)
	}
	rng := rand.New(rand.NewSource(seed))

	_, _ = ntcRijqFUuoGg.WriteString("\n\n        <h1>Hello ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcbGHNBxxTTK.EscapeString(name))
	_, _ = ntcRijqFUuoGg.WriteString("!<h1>\n        <p>Here's your lucky number: ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcVqcLceiEsr.FormatInt(int64(rng.Intn(100)), 10))
	_, _ = ntcRijqFUuoGg.WriteString("</p>\n\n    </body>\n    </html>\n")
	return ntcRijqFUuoGg.String()
}
