package nc

import (
	"sync/atomic"
)

// Global garbageman.
var garbageman Garbageman

// Get a buffer from the global garbageman.
func Buffer() Block {
	b := MakeBlock(WarpSize())
	b.refcount = new(Refcount)
	return b
	//return garbageman.Get()
}

func Buffer3() [3]Block {
	return [3]Block{Buffer(), Buffer(), Buffer()}
}

// Recyle a buffer from the global garbageman.
func Recycle(g ...Block) {
	//garbageman.Recycle(g...) 
}

func Recycle3(garbages ...[3]Block) {
	for _, g := range garbages {
		for _, c := range g {
			garbageman.Recycle(c)
		}
	}
}

func NumAlloc() int { return int(atomic.LoadInt32(&garbageman.numAlloc)) }

// Global garbageman.
var gpugarbageman GpuGarbageman

// Get a buffer form the global garbageman.
func GpuBuffer() GpuBlock {
	return GpuBufferSize(WarpSize())
}

func GpuBufferSize(size [3]int) GpuBlock {
	//return gpugarbageman.GetSize(size)
	b := MakeGpuBlock(size)
	b.refcount = new(Refcount)
	return b
}

func RecycleGpu(garbages ...GpuBlock) {
	for _, g := range garbages {
		if g.refcount == nil {
			continue // slice does not originate from here
		}
		if g.refcount.Load() == 0 {
			g.Free()
		} else { // cannot be recycled, just yet
			g.refcount.Add(-1)
		}
	}
	//gpugarbageman.Recycle(g...)
}

func NumGpuAlloc() int {
	return int(atomic.LoadInt32(&gpugarbageman.numAlloc))
}

func InitGarbageman() {
	// recycling buffer may be huge, it should not waste any memory.
	garbageman.Init(WarpSize())
	gpugarbageman.Init(WarpSize())
}

// Garbage chute buffer size (blocks)
const BUFSIZE = 1000