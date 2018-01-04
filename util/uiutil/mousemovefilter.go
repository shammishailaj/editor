package uiutil

import (
	"time"

	"github.com/jmigpin/editor/util/uiutil/event"
)

func MouseMoveFilterLoop(in <-chan interface{}, out chan<- interface{}, fps *int) {
	var lastMoveEv interface{}
	var ticker *time.Ticker
	var timeToSend <-chan time.Time

	//n := 0
	keepMoveEv := func(ev interface{}) {
		//n++
		lastMoveEv = ev
		if ticker == nil {
			ticker = time.NewTicker(time.Second / time.Duration(*fps))
			timeToSend = ticker.C
		}
	}

	sendMoveEv := func() {
		//log.Printf("kept %d times before sending", n)
		//n = 0
		ticker.Stop()
		ticker = nil
		timeToSend = nil
		out <- lastMoveEv
	}

	sendMoveEvIfKept := func() {
		if ticker != nil {
			sendMoveEv()
		}
	}

	for {
		select {
		case ev, ok := <-in:
			if !ok {
				goto forEnd
			}

			isMove := false
			if wi, ok := ev.(*event.WindowInput); ok {
				if _, ok := wi.Event.(*event.MouseMove); ok {
					isMove = true
					keepMoveEv(ev)
				}
			}
			if !isMove {
				sendMoveEvIfKept()
				out <- ev
			}
		case <-timeToSend:
			sendMoveEv()
		}
	}
forEnd:
}