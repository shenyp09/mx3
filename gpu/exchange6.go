package gpu

import (
	"github.com/barnex/cuda5/cu"
	"github.com/barnex/cuda5/safe"
	"nimble-cube/core"
	"nimble-cube/gpu/ptx"
	"unsafe"
)

type Exchange6 struct {
	m   RChan3
	hex Chan3
	*core.Mesh
	aex_reduced float64
	factors     [3]float32
	stream      cu.Stream
}

func NewExchange6(m RChan3, hex Chan3, mesh *core.Mesh, aex_reduced float64) *Exchange6 {
	e := &Exchange6{m: m, hex: hex, Mesh: mesh, aex_reduced: aex_reduced, stream: cu.StreamCreate()}
	cellsize := mesh.CellSize()
	for i := range e.factors {
		e.factors[i] = float32(aex_reduced / (cellsize[i] * cellsize[i]))
	}
	return e
}

func (e *Exchange6) Run() {
	LockCudaThread()
	N := e.NCell()
	for {
		m := e.m.ReadNext(N)
		hex := e.hex.WriteNext(N)
		for i := 0; i < 3; i++ {
			exchange6(hex[i], m[i], e.Mesh, e.factors, e.stream)
		}
		e.hex.WriteDone()
		e.m.ReadDone()
	}
}

var exchange6Code cu.Function

//__global__ void exchange6(float* h, float* m, float fac0, float fac1, float fac2, int wrap0, int wrap1, int wrap2, int N0, int N1, int N2){
func exchange6(h, m safe.Float32s, mesh *core.Mesh, factors [3]float32, stream cu.Stream) {
	core.Assert(h.Len() == m.Len() && h.Len() == mesh.NCell())

	if exchange6Code == 0 {
		mod := cu.ModuleLoadData(ptx.EXCHANGE6)
		exchange6Code = mod.GetFunction("exchange6")
	}

	size := mesh.GridSize()
	N0, N1, N2 := size[0], size[1], size[2]
	wrap := mesh.PBC()

	gridDim, blockDim := Make2DConf(N2, N1) // why?

	hptr := h.Pointer()
	mptr := m.Pointer()

	args := []unsafe.Pointer{
		unsafe.Pointer(&hptr),
		unsafe.Pointer(&mptr),
		unsafe.Pointer(&factors[0]),
		unsafe.Pointer(&factors[1]),
		unsafe.Pointer(&factors[2]),
		unsafe.Pointer(&wrap[0]),
		unsafe.Pointer(&wrap[1]),
		unsafe.Pointer(&wrap[2]),
		unsafe.Pointer(&N0),
		unsafe.Pointer(&N1),
		unsafe.Pointer(&N2)}

	shmem := 0
	cu.LaunchKernel(exchange6Code, gridDim.X, gridDim.Y, gridDim.Z, blockDim.X, blockDim.Y, blockDim.Z, shmem, stream, args)
	stream.Synchronize()
}