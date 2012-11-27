package gpu

import (
	"code.google.com/p/nimble-cube/core"
	"code.google.com/p/nimble-cube/nimble"
	"github.com/barnex/cuda5/cu"
	"github.com/barnex/cuda5/safe"
	"unsafe"
)

type Symm2D struct {
	size        [3]int              // 3D size of the input/output data
	kernSize    [3]int              // Size of kernel and logical FFT size.
	fftKernSize [3]int              // Size of real, FFTed kernel
	n           int                 // product of size
	input       [3]nimble.RChan1    // TODO: fuse with input
	output      [3]nimble.Chan1     // TODO: fuse with output
	outChan     nimble.ChanN        // same as output
	fftRBuf     [3]safe.Float32s    // FFT input buf for FFT, shares storage with fftCBuf. 
	fftCBuf     [3]safe.Complex64s  // FFT output buf, shares storage with fftRBuf
	gpuFFTKern  [3][3]safe.Float32s // FFT kernel on device: TODO: xfer if needed
	fwPlan      safe.FFT3DR2CPlan   // Forward FFT (1 component)
	bwPlan      safe.FFT3DC2RPlan   // Backward FFT (1 component)
	stream      cu.Stream           // 
	kern        [3][3][]float32     // Real-space kernel
	kernArr     [3][3][][][]float32 // Real-space kernel
	inited      bool
}

func (c *Symm2D) init() {
	core.Log("initializing 2D symmetric convolution")
	if c.inited {
		core.Panic("conv: already initialized")
	}
	c.inited = true
	// TODO: may unlock main thread...
	LockCudaThread()
	defer UnlockCudaThread()

	padded := c.kernSize

	{ // init FFT plans
		c.stream = cu.StreamCreate()
		c.fwPlan = safe.FFT3DR2C(padded[0], padded[1], padded[2])
		c.fwPlan.SetStream(c.stream)
		c.bwPlan = safe.FFT3DC2R(padded[0], padded[1], padded[2])
		c.bwPlan.SetStream(c.stream)
	}

	{ // init device buffers
		// 2D re-uses fftBuf[1] as fftBuf[0], 3D needs all 3 fftBufs.
		for i := 1; i < 3; i++ {
			c.fftCBuf[i] = MakeComplexs(prod(fftR2COutputSizeFloats(c.kernSize)) / 2)
		}
		if c.is3D() {
			c.fftCBuf[0] = MakeComplexs(prod(fftR2COutputSizeFloats(c.kernSize)) / 2)
		} else {
			c.fftCBuf[0] = c.fftCBuf[1]
		}
		for i := 0; i < 3; i++ {
			c.fftRBuf[i] = c.fftCBuf[i].Float().Slice(0, prod(c.kernSize))
		}
	}

	if c.is2D() {
		c.initFFTKern2D()
	} else {
		c.initFFTKern3D()
	}
}

func (c *Symm2D) initFFTKern3D() {
	padded := c.kernSize
	ffted := fftR2COutputSizeFloats(padded)
	realsize := ffted
	realsize[2] /= 2
	c.fftKernSize = realsize
	halfkern := realsize
	//halfkern[1] = halfkern[1]/2 + 1
	fwPlan := c.fwPlan
	output := safe.MakeComplex64s(fwPlan.OutputLen())
	defer output.Free()
	input := output.Float().Slice(0, fwPlan.InputLen())

	// upper triangular part
	fftKern := make([]float32, prod(halfkern))
	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			if c.kern[i][j] != nil { // ignore 0's
				input.CopyHtoD(c.kern[i][j])
				fwPlan.Exec(input, output)
				fwPlan.Stream().Synchronize() // !!
				scaleRealParts(fftKern, output.Float().Slice(0, prod(halfkern)*2), 1/float32(fwPlan.InputLen()))
				c.gpuFFTKern[i][j] = safe.MakeFloat32s(len(fftKern))
				c.gpuFFTKern[i][j].CopyHtoD(fftKern)
			}
		}
	}
}

// Initialize GPU FFT kernel for 2D. 
// Only the non-redundant parts are stored on the GPU.
func (c *Symm2D) initFFTKern2D() {
	padded := c.kernSize
	ffted := fftR2COutputSizeFloats(padded)
	realsize := ffted
	realsize[2] /= 2
	c.fftKernSize = realsize
	halfkern := realsize
	halfkern[1] = halfkern[1]/2 + 1
	fwPlan := c.fwPlan
	output := HostFloats(2 * fwPlan.OutputLen()).Complex()
	defer cu.MemFreeHost(unsafe.Pointer(uintptr(output.Pointer()))) // TODO: is Float32s safe with uintptr?
	input := output.Float().Slice(0, fwPlan.InputLen())

	// upper triangular part
	fftKern := make([]float32, prod(halfkern))
	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ {
			if c.kern[i][j] != nil { // ignore 0's
				input.CopyHtoD(c.kern[i][j])
				fwPlan.Exec(input, output)
				fwPlan.Stream().Synchronize() // !!
				scaleRealParts(fftKern, output.Float().Slice(0, prod(halfkern)*2), 1/float32(fwPlan.InputLen()))
				c.gpuFFTKern[i][j] = MakeFloats(len(fftKern))
				c.gpuFFTKern[i][j].CopyHtoD(fftKern)
			}
		}
	}
}

func (c *Symm2D) Run() {
	core.Log("running symmetric 2D convolution")
	LockCudaThread()

	for {
		c.Exec()
	}
}

func (c *Symm2D) Exec() {
	if c.is2D() {
		c.exec2D()
	} else {
		c.exec3D()
	}
}

func (c *Symm2D) exec3D() {
	padded := c.kernSize
	offset := [3]int{0, 0, 0}

	//N0, N1, N2 := cc.fftKernSize[1], c.fftKernSize[2]
	for i := 0; i < 3; i++ {
		in := c.input[i].ReadNext(c.n).Device()
		c.fftRBuf[i].MemsetAsync(0, c.stream)
		copyPad(c.fftRBuf[i], in, padded, c.size, offset, c.stream)
		c.fwPlan.Exec(c.fftRBuf[i], c.fftCBuf[i])
		c.stream.Synchronize()
		c.input[i].ReadDone()
	}

	// kern mul
	kernMulRSymm(c.fftCBuf,
		c.gpuFFTKern[0][0], c.gpuFFTKern[1][1], c.gpuFFTKern[2][2],
		c.gpuFFTKern[1][2], c.gpuFFTKern[0][2], c.gpuFFTKern[0][1],
		c.stream)
	c.stream.Synchronize()

	// BW FFT 
	for i := 0; i < 3; i++ {
		out := c.output[i].WriteNext(c.n).Device()
		c.bwPlan.Exec(c.fftCBuf[i], c.fftRBuf[i])
		copyPad(out, c.fftRBuf[i], c.size, padded, offset, c.stream)
		c.stream.Synchronize()
		c.output[i].WriteDone()
	}
}

func (c *Symm2D) exec2D() {
	padded := c.kernSize
	offset := [3]int{0, 0, 0}

	N1, N2 := c.fftKernSize[1], c.fftKernSize[2]
	// Convolution is separated into 
	// a 1D convolution for x
	// and a 2D convolution for yz.
	// so only 2 FFT buffers are then needed at the same time.

	// FFT x
	in := c.input[0].ReadNext(c.n).Device()
	c.fftRBuf[0].MemsetAsync(0, c.stream) // copypad does NOT zero remainder.
	copyPad(c.fftRBuf[0], in, padded, c.size, offset, c.stream)
	c.fwPlan.Exec(c.fftRBuf[0], c.fftCBuf[0])
	//c.stream.Synchronize()
	c.input[0].ReadDone()

	// kern mul X	
	kernMulRSymm2Dx(c.fftCBuf[0], c.gpuFFTKern[0][0], N1, N2, c.stream)
	//c.stream.Synchronize()

	// bw FFT x
	out := c.output[0].WriteNext(c.n).Device()
	c.bwPlan.Exec(c.fftCBuf[0], c.fftRBuf[0])
	copyPad(out, c.fftRBuf[0], c.size, padded, offset, c.stream)
	c.stream.Synchronize()
	c.output[0].WriteDone()

	// FW FFT yz
	for i := 1; i < 3; i++ {
		in := c.input[i].ReadNext(c.n).Device()
		c.fftRBuf[i].MemsetAsync(0, c.stream)
		copyPad(c.fftRBuf[i], in, padded, c.size, offset, c.stream)
		c.fwPlan.Exec(c.fftRBuf[i], c.fftCBuf[i])
		c.stream.Synchronize()
		c.input[i].ReadDone()
	}

	// kern mul yz
	kernMulRSymm2Dyz(c.fftCBuf[1], c.fftCBuf[2],
		c.gpuFFTKern[1][1], c.gpuFFTKern[2][2], c.gpuFFTKern[1][2],
		N1, N2, c.stream)
	c.stream.Synchronize()

	// BW FFT yz
	for i := 1; i < 3; i++ {
		out := c.output[i].WriteNext(c.n).Device()
		c.bwPlan.Exec(c.fftCBuf[i], c.fftRBuf[i])
		copyPad(out, c.fftRBuf[i], c.size, padded, offset, c.stream)
		c.stream.Synchronize()
		c.output[i].WriteDone()
	}
}

func (c *Symm2D) is2D() bool {
	return c.size[0] == 1
}

func (c *Symm2D) is3D() bool {
	return !c.is2D()
}

func (c *Symm2D) Output() nimble.ChanN {
	return c.outChan
}

func NewConvolution(tag, unit string, mesh *nimble.Mesh, memType nimble.MemType, kernel [3][3][][][]float32, input_ nimble.ChanN) *Symm2D {
	size := mesh.Size()
	in_ := input_.NewReader()
	input := [3]nimble.RChan1{in_[0], in_[1], in_[2]}
	c := new(Symm2D)
	c.size = size
	c.kernArr = kernel
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if kernel[i][j] != nil {
				c.kern[i][j] = core.Contiguous(kernel[i][j])
			}
		}
	}
	c.n = prod(size)
	c.kernSize = core.SizeOf(kernel[0][0])
	c.input = input
	c.outChan = nimble.MakeChanN(3, tag, unit, mesh, memType, 0)
	c.output = [3]nimble.Chan1{c.outChan.Comp(0), c.outChan.Comp(1), c.outChan.Comp(2)}

	c.init()
	nimble.Stack(c)

	return c
	// TODO: self-test
}