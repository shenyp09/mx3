package data

import (
	"code.google.com/p/mx3/util"
	"hash"
	"hash/crc64"
	"io"
	"math"
	"os"
	"unsafe"
)

// TODO: buffer? benchmark
func Write(out io.Writer, s *Slice, time float64) error {
	w := newWriter(out)
	w.header.Components = s.NComp()
	w.header.MeshSize = s.Mesh().Size()
	w.header.MeshStep = s.Mesh().CellSize()
	w.header.MeshUnit = "m"
	w.header.Time = time
	w.header.TimeUnit = "s"
	w.header.DataLabel = s.Tag()
	w.header.DataUnit = s.Unit()
	w.header.Precission = 4
	w.writeHeader()
	if w.err != nil {
		return w.err
	}

	data := s.Host()
	for _, d := range data {
		w.writeData(d)
	}
	w.writeHash()
	return w.err
}

func WriteFile(fname string, s *Slice, time float64) error {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return Write(f, s, time)
}

func MustWriteFile(fname string, s *Slice, time float64) {
	err := WriteFile(fname, s, time)
	util.FatalErr(err)
}

var table = crc64.MakeTable(crc64.ISO)

type writer struct {
	header
	out io.Writer
	crc hash.Hash64
	err error
}

func newWriter(out io.Writer) *writer {
	w := new(writer)
	w.crc = crc64.New(table)
	w.out = io.MultiWriter(w.crc, out)
	return w
}

const MAGIC = "#dump002"
const padding = 0

// Writes the current header.
func (w *writer) writeHeader() {
	w.writeString(MAGIC)
	w.writeUInt64(uint64(w.Components))
	for _, s := range w.MeshSize {
		w.writeUInt64(uint64(s))
	}
	for _, s := range w.MeshStep {
		w.writeFloat64(s)
	}
	w.writeString(w.MeshUnit)
	w.writeFloat64(w.Time)
	w.writeString(w.TimeUnit)
	w.writeString(w.DataLabel)
	w.writeString(w.DataUnit)
	w.writeUInt64(w.Precission)

	for i := 0; i < padding; i++ {
		w.writeUInt64(0)
	}
}

// Writes the data.
func (w *writer) writeData(list []float32) {
	w.count(w.out.Write((*(*[1<<31 - 1]byte)(unsafe.Pointer(&list[0])))[0 : 4*len(list)]))
}

// Writes the accumulated hash of this frame, closing the frame.
func (w *writer) writeHash() {
	w.writeUInt64(w.crc.Sum64())
	w.crc.Reset()
}

func (w *writer) count(n int, err error) {
	if err != nil && w.err == nil {
		w.err = err
	}
}

func (w *writer) writeFloat64(x float64) {
	w.writeUInt64(math.Float64bits(x))
}

func (w *writer) writeString(x string) {
	var buf [8]byte
	copy(buf[:], x)
	w.count(w.out.Write(buf[:]))
}

func (w *writer) writeUInt64(x uint64) {
	w.count(w.out.Write((*(*[8]byte)(unsafe.Pointer(&x)))[:8]))
}
