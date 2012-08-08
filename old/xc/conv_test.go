package xc

import (
	"nimble-cube/core"
	"testing"
)

func TestConv2(test *testing.T) {
	size := [3]int{1, 32, 32}
	core.InitSize(size[0], size[1], size[2])
	core.InitCellSize(1e-9, 1e-9, 1e-9)
	N := prod(size)

	in := make([]float32, 3*N)
	input := [3][]float32{in[0*N : 1*N], in[1*N : 2*N], in[2*N : 3*N]}

	out := make([]float32, 3*N)
	output := [3][]float32{out[0*N : 1*N], out[1*N : 2*N], out[2*N : 3*N]}

	conv := NewConv2(input, output, size)

	TestConv(conv)
}

func BenchmarkConv1_2DSmall(b *testing.B) {
	b.StopTimer()

	size := [3]int{1, 128, 128}
	core.InitSize(size[0], size[1], size[2])
	core.InitCellSize(1e-9, 1e-9, 1e-9)
	N := prod(size)

	in := make([]float32, 3*N)
	input := [3][]float32{in[0*N : 1*N], in[1*N : 2*N], in[2*N : 3*N]}

	out := make([]float32, 3*N)
	output := [3][]float32{out[0*N : 1*N], out[1*N : 2*N], out[2*N : 3*N]}

	conv := NewConv1(input, output, size)

	b.SetBytes(int64(prod(size)) * 4 * 2) // *2: xfer back and forth

	conv.Test()

	core.DEBUG = false
	core.LOG = false

	// warmup
	conv.Push(core.N())
	conv.Pull(core.N())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		conv.Push(core.N())
		conv.Pull(core.N())
	}
	b.StopTimer()
	core.Cleanup()
}