package parse

import "fmt"

type FileSet struct {
	Counter   int64
	Filenames map[string]int64
	Newlines  map[string][]int64
}

func NewFileSet() *FileSet {
	return &FileSet{
		Filenames: make(map[string]int64),
		Newlines:  make(map[string][]int64),
	}
}

func (fs *FileSet) AddFile(filename string, contents []byte) int64 {
	fs.Filenames[filename] = fs.Counter
	var x []int64
	for i, b := range contents {
		if b == '\n' {
			x = append(x, fs.Counter+int64(i))
		}
	}
	fs.Newlines[filename] = x
	bic := fs.Counter
	fs.Counter += int64(len(contents))
	return bic
}

func (fs *FileSet) ResolvePosition(pos int64) string {
	var filename string
	for fname, fpos := range fs.Filenames {
		if pos < fpos {
			break
		}
		if pos > fpos {
			filename = fname
		}
	}

	var linepos int64
	for _, lineStart := range fs.Newlines[filename] {
		if pos < lineStart {
			break
		}
		if pos > lineStart {
			linepos = lineStart
		}
	}

	line := linepos - fs.Filenames[filename]
	col := pos - linepos
	return fmt.Sprintf("%s:%d:%d", filename, line, col)
}
