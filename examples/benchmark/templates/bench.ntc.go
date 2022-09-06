package templates

import (
	ntcTxhiXmdYod "html"
	ntcBAMGdhfBWH "strconv"
	ntcSuZXQQTKkx "strings"
)

type BenchRow struct {
	ID      int
	Message string
	Print   bool
}

func BenchPage(rows []BenchRow) string {
	ntcRijqFUuoGg := new(ntcSuZXQQTKkx.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("<html>\n    <head><title>test</title></head>\n    \t<body>\n    \t\t<ul>\n    \t\t")
	for _, row := range rows {
		_, _ = ntcRijqFUuoGg.WriteString("\n    \t\t\t")
		if row.Print {
			_, _ = ntcRijqFUuoGg.WriteString("\n    \t\t\t\t<li>ID=")
			_, _ = ntcRijqFUuoGg.WriteString(ntcBAMGdhfBWH.FormatInt(int64(row.ID), 10))
			_, _ = ntcRijqFUuoGg.WriteString(", Message=")
			_, _ = ntcRijqFUuoGg.WriteString(ntcTxhiXmdYod.EscapeString(row.Message))
			_, _ = ntcRijqFUuoGg.WriteString("</li>\n    \t\t\t")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n    \t\t")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n    \t\t</ul>\n    \t</body>\n    </html>\n")
	return ntcRijqFUuoGg.String()
}
