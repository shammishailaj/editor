package ui

import "github.com/jmigpin/editor/xutil/xgbutil"

type Toolbar struct {
	*TextArea
	OnSetText func()
}

func NewToolbar(ui *UI) *Toolbar {
	tb := &Toolbar{TextArea: NewTextArea(ui)}
	tb.DisableHighlightCursorWord = true
	tb.DisablePageUpDown = true

	tb.TextArea.EvReg.Add(TextAreaSetStrEventId,
		&xgbutil.ERCallback{tb.onTextAreaSetStr})

	return tb
}
func (tb *Toolbar) onTextAreaSetStr(ev0 xgbutil.EREvent) {
	ev := ev0.(*TextAreaSetStrEvent)
	if tb.OnSetText != nil {
		tb.OnSetText()
	}
	// keep pointer inside if it was in before
	// useful in dynamic bounds becoming shorter and leaving the pointer outside, losing keyboard focus
	p, ok := tb.ui.Win.QueryPointer()
	if ok && p.In(ev.OldBounds) && !p.In(tb.C.Bounds) {
		tb.ui.WarpPointerToRectanglePad(&tb.C.Bounds)
	}
}
