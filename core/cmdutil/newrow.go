package cmdutil

import (
	"log"
	"os"
	"path"
)

func NewRow(ed Editorer) {
	p, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}

	col, nextRow := ed.GoodColumnRowPlace()

	erow2, ok := ed.ActiveERow()
	if ok {
		fp := erow2.Filename()

		// stick with directory if exists, otherwise get base dir
		fi, err := os.Stat(fp)
		if err == nil && fi.IsDir() {
			p = fp
		} else {
			p = path.Dir(fp)
		}

		// position after active row
		r := erow2.Row()
		col = r.Col
		nr, ok := r.NextRow()
		if ok {
			nextRow = nr
		}
	}

	erow := ed.NewERowBeforeRow(p, col, nextRow)
	erow.Row().WarpPointer()
}
