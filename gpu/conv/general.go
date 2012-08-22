package conv

import (
	"github.com/barnex/cuda4/cu"
	"github.com/barnex/cuda4/safe"
	"nimble-cube/core"
)

// General convolution, not optimized for specific cases.
// Also not concurrent.
type General struct {
	hostData    // sizes, host input/output/kernel arrays
	deviceData3 // device buffers // could use just one ioBuf
	fwPlan      safe.FFT3DR2CPlan
	bwPlan      safe.FFT3DC2RPlan
}

func (c *General) Exec() {
	for i := 0; i < 3; i++ {
		core.Print("input", i, "\n", core.Reshape(c.input[i], c.IOSize()))
		c.ioBuf[i].CopyHtoD(c.input[i])
		core.Print("ioBuf", i, "\n", core.Reshape(c.ioBuf[i].Host(), c.IOSize()))
		c.copyPadIOBuf(i)
		core.Print("fftRBuf", i, "\n", core.Reshape(c.fftRBuf[i].Host(), c.KernelSize()))
	}
}

// Copy ioBuf[i] to fftRBuf[i], adding padding :-).
func (c *General) copyPadIOBuf(i int) {
	stream0 := cu.Stream(0)
	offset := [3]int{0, 0, 0}
	c.fftRBuf[i].Memset(0) // copypad does NOT zero remainder.
	stream0.Synchronize()
	copyPad(c.fftRBuf[i], c.ioBuf[i], c.kernSize, c.size, offset, stream0)
	stream0.Synchronize()
	c.fwPlan.Exec(c.fftRBuf[i], c.fftCBuf[i])
}

// Size of the FFT'ed kernel expressed in number of floats.
func (c *General) FFTKernelSizeFloats() [3]int {
	return fftR2COutputSizeFloats(c.KernelSize())
	// kernel size is FFT logic size
}

// Initializes c.gpuFFTKern and c.fftKern
func (c *General) initFFTKern() {
	realsize := c.kernSize
	reallen := prod(realsize)
	fftedsize := fftR2COutputSizeFloats(realsize)
	fftedlen := prod(fftedsize)

	fwPlan := c.fwPlan // could use any

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			c.fftKern[i][j] = make([]float32, fftedlen)
			c.gpuFFTKern[i][j] = safe.MakeFloat32s(fftedlen)
			c.gpuFFTKern[i][j].Slice(0, reallen).CopyHtoD(c.kern[i][j])
			fwPlan.Exec(c.gpuFFTKern[i][j].Slice(0, reallen), c.gpuFFTKern[i][j].Complex())
			c.gpuFFTKern[i][j].CopyDtoH(c.fftKern[i][j])
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			core.Print("kern", i, j, "\n", core.Reshape(c.kern[i][j], c.kernSize))
			core.Print("fftKern", i, j, "\n", core.Reshape(c.fftKern[i][j], fftedsize))
		}
	}

}

// Initializes the FFT plans.
func (c *General) initFFT() {
	padded := c.kernSize
	//realsize := fftR2COutputSizeFloats(padded)
	c.fwPlan = safe.FFT3DR2C(padded[0], padded[1], padded[2])
	c.bwPlan = safe.FFT3DC2R(padded[0], padded[1], padded[2])
	// no streams set yet
}

func NewGeneral(input_, output_ [3][][][]float32, kernel [3][3][][][]float32) *General {
	c := new(General)
	c.hostData.init(input_, output_, kernel)

	// need cuda thread lock from here on:
	c.hostData.initPageLock()
	c.initFFT()
	c.initFFTKern()
	c.deviceData3.init(c.IOSize(), c.KernelSize())

	return c
}