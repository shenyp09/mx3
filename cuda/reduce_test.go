package cuda

import (
	"code.google.com/p/mx3/data"
	"testing"
)

// test input data
var in1, in2, in3 *data.Slice

func initTest() {
	if in1 != nil {
		return
	}
	LockThread()
	Init()
	{
		inh1 := make([]float32, 1000)
		for i := range inh1 {
			inh1[i] = float32(i)
		}
		in1 = toGPU(inh1)
	}
	{
		inh2 := make([]float32, 100000)
		for i := range inh2 {
			inh2[i] = -float32(i) / 100
		}
		in2 = toGPU(inh2)
	}
}

func toGPU(list []float32) *data.Slice {
	mesh := data.NewMesh(1, 1, len(list), 1, 1, 1)
	h := data.SliceFromList([][]float32{list}, mesh)
	d := NewSlice(1, mesh)
	data.Copy(d, h)
	return d
}

func TestReduceSum(t *testing.T) {
	initTest()
	result := Sum(in1)
	if result != 499500 {
		t.Error("got:", result)
	}
}

//func TestReduceMax(t *testing.T) {
//	initTest()
//	result := Max(in1)
//	if result != 999 {
//		t.Error("got:", result)
//	}
//	result = Max(in2)
//	if result != 0 {
//		t.Error("got:", result)
//	}
//}
//
//func TestReduceMin(t *testing.T) {
//	result := Min(in1)
//	if result != 0 {
//		t.Error("got:", result)
//	}
//
//	result = Min(in2)
//	if result != -999.99 {
//		t.Error("got:", result)
//	}
//}

func TestReduceMaxAbs(t *testing.T) {
	result := MaxAbs(in1)
	if result != 999 {
		t.Error("got:", result)
	}
	result = MaxAbs(in2)
	if result != 999.99 {
		t.Error("got:", result)
	}
}

//func TestReduceMaxDiff(t *testing.T) {
//	LockCudaThread()
//	N := 100001
//	input := nimble.MakeSlice(N, nimble.UnifiedMemory)
//	in := input.Host()
//	for i := range in {
//		in[i] = -float32(i) / 100
//	}
//	input2 := nimble.MakeSlice(N, nimble.UnifiedMemory)
//	in2 := input2.Host()
//	for i := range in2 {
//		in2[i] = float32(i) / 100
//	}
//	result := MaxDiff(input.Device(), input2.Device())
//	if result != 2000 {
//		t.Error("got:", result)
//	}
//}
//
//func BenchmarkReduceSum(b *testing.B) {
//	core.LOG = false
//	b.StopTimer()
//	LockCudaThread()
//	const N = 32 * 1024 * 1024
//	input := nimble.MakeSlice(N, nimble.GPUMemory)
//	b.SetBytes(N * 4)
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		Sum(input.Device())
//	}
//}
//
//func TestReduceMaxVecNorm(t *testing.T) {
//	LockCudaThread()
//	N := 1234
//	input := nimble.MakeSlice(N, nimble.UnifiedMemory)
//	in := input.Host()
//	for i := range in {
//		in[i] = -float32(i) / 1000
//	}
//	x := input.Device()
//	result := MaxVecNorm(x, x, x)
//	want := math.Sqrt(3) * 1233. / 1000.
//	if math.Abs(result-want) > 1e-7 {
//		t.Error("got:", result, "want:", want)
//	}
//}
//
//func TestReduceMaxVecDiff(t *testing.T) {
//	LockCudaThread()
//	N := 1234
//	input := nimble.MakeSlice(N, nimble.UnifiedMemory)
//	in := input.Host()
//	for i := range in {
//		in[i] = -float32(i) / 1000
//	}
//	x := input.Device()
//	input2 := nimble.MakeSlice(N, nimble.UnifiedMemory)
//	in2 := input2.Host()
//	for i := range in2 {
//		in2[i] = 0
//	}
//	y := input2.Device()
//	result := MaxVecDiff(x, x, x, y, y, y)
//	want := math.Sqrt(3) * 1233. / 1000.
//	if math.Abs(result-want) > 1e-7 {
//		t.Error("got:", result, "want:", want)
//	}
//}
