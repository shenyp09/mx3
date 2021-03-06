package engine

import (
	"code.google.com/p/mx3/cuda"
	"code.google.com/p/mx3/data"
	"code.google.com/p/mx3/util"
	"log"
)

// function that adds a quantity to dst
type addFunc func(dst *data.Slice)

// Output Handle for a quantity that is not explicitly stored,
// but only added to an other quantity (like effective field)
type adder struct {
	addFn addFunc // calculates quantity and add result to dst
	autosave
}

func newAdder(name string, f addFunc) *adder {
	a := new(adder)
	a.addFn = f
	a.name = name
	return a
}

// Calls the addFunc to add the quantity to Dst. If output is needed,
// it is first added to a separate buffer, saved, and then added to Dst.
func (a *adder) addTo(Dst *buffered, goodstep bool) {
	if goodstep && a.needSave() {
		Buf := addBuf()
		buf := Buf.Write()
		cuda.Zero(buf)
		a.addFn(buf)
		dst := Dst.Write()
		cuda.Madd2(dst, dst, buf, 1, 1)
		Dst.WriteDone()
		// Buf only unlocked here to avoid overwite by next addTo
		goSave(a.fname(), buf, Time, func() { Buf.WriteDone() })
		a.saved()
	} else {
		dst := Dst.Write()
		a.addFn(dst)
		Dst.WriteDone()
	}
}

var _addBuf *buffered

// returns a GPU buffer for temporarily adding a quantity to and saving it
func addBuf() *buffered {
	if _addBuf == nil {
		util.DashExit()
		log.Println("allocating GPU buffer for output")
		_addBuf = newBuffered(cuda.NewSynced(3, mesh), "buffer", nil)
	}
	return _addBuf
}
