package textutil

import (
	"image"
	"io"

	"github.com/jmigpin/editor/util/iout"
	"github.com/jmigpin/editor/util/uiutil/event"
	"github.com/jmigpin/editor/util/uiutil/widget"
)

func MoveCursorToPoint(te *widget.TextEdit, p *image.Point, sel bool) {
	tc := te.TextCursor

	i := te.GetIndex(*p)
	tc.SetSelectionUpdate(sel, i)

	// set primary copy
	if tc.SelectionOn() {
		s, err := tc.Selection()
		if err == nil {
			te.SetCPCopy(event.CPIPrimary, string(s))
		}
	}
}

//----------

func MoveCursorLeft(te *widget.TextEdit, sel bool) error {
	tc := te.TextCursor
	ci := tc.Index()
	_, size, err := tc.RW().ReadLastRuneAt(ci)
	if err != nil {
		return err
	}
	tc.SetSelectionUpdate(sel, ci-size)
	return nil
}

func MoveCursorRight(te *widget.TextEdit, sel bool) error {
	tc := te.TextCursor
	ci := tc.Index()
	_, size, err := tc.RW().ReadRuneAt(ci)
	if err != nil {
		return err
	}
	tc.SetSelectionUpdate(sel, ci+size)
	return nil
}

//----------

func MoveCursorUp(te *widget.TextEdit, sel bool) {
	tc := te.TextCursor

	p := te.GetPoint(tc.Index())
	p.Y -= te.LineHeight() - 1
	i := te.GetIndex(p)

	tc.SetSelectionUpdate(sel, i)
}

func MoveCursorDown(es *widget.TextEdit, sel bool) {
	tc := es.TextCursor

	p := es.GetPoint(tc.Index())
	p.Y += es.LineHeight() + 1
	i := es.GetIndex(p)

	tc.SetSelectionUpdate(sel, i)
}

//----------

func MoveCursorJumpLeft(te *widget.TextEdit, sel bool) error {
	tc := te.TextCursor
	i, err := jumpLeftIndex(tc)
	if err != nil {
		return err
	}
	tc.SetSelectionUpdate(sel, i)
	return nil
}
func MoveCursorJumpRight(te *widget.TextEdit, sel bool) error {
	tc := te.TextCursor
	i, err := jumpRightIndex(tc)
	if err != nil {
		return err
	}
	tc.SetSelectionUpdate(sel, i)
	return nil
}

//----------

func jumpLeftIndex(tc *widget.TextCursor) (int, error) {
	i, size, err := iout.LastIndexFunc(tc.RW(), tc.Index(), 1000, true, edgeOfNextWordOrNewline())
	if err != nil {
		if err == io.EOF {
			i = 0
		} else {
			return 0, err
		}
	}
	return i + size, nil
}

func jumpRightIndex(tc *widget.TextCursor) (int, error) {
	i, _, err := iout.IndexFunc(tc.RW(), tc.Index(), 1000, true, edgeOfNextWordOrNewline())
	if err != nil {
		if err == io.EOF {
			i = tc.RW().Len()
		} else {
			return 0, err
		}
	}
	return i, nil
}

//----------

func edgeOfNextWordOrNewline() func(rune) bool {
	first := true
	var inWord bool
	return func(ru rune) bool {
		w := isWordRune(ru)
		if first {
			first = false
			inWord = w
		} else {
			if !inWord {
				inWord = w

				if ru == '\n' {
					return true
				}
			} else {
				if !w {
					return true
				}
			}
		}
		return false
	}
}
