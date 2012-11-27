package gpu

import (
	"code.google.com/p/nimble-cube/nimble"
	"github.com/barnex/cuda5/cu"
	"testing"
)

func TestReduceSum(t *testing.T) {
	LockCudaThread()
	N := 512
	input := nimble.MakeSlice(N, nimble.UnifiedMemory)
//		in := input.Host()
//		for i := range in {
//			in[i] = 0
//		}
	str := cu.StreamCreate()
	result := reduce_sum(input.Device(), str)
	if result != 0 {
		t.Error("got:", result)
	}
}