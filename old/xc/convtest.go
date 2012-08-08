package xc

import (
	"fmt"
	"github.com/barnex/cuda4/safe"
	"github.com/barnex/fmath"
	"math"
	"math/rand"
	"nimble-cube/core"
)

func TestConv(conv *Conv2) {
	core.Log("Testing convolution")
	size := conv.size
	N := prod(size)

	in_ := conv.input
	var in [3][][][]float32
	for c := 0; c < 3; c++ {
		in[c] = safe.Reshape3DFloat32(in_[c], size[0], size[1], size[2])
	}

	for i := range in_[0] {
		in_[0][i] = 1
	}

	ref := BruteSymmetricConvolution(conv.input, conv.kern, size)

	conv.Push(N)
	conv.Pull(N)

	core.Log("brute:", ref)
	core.Log("ffted:", conv.output)

	var max float32
	for c := 0; c < 3; c++ {
		for i := range ref[c] {
			if fmath.Abs(ref[c][i]-conv.output[c][i]) > max {
				max = fmath.Abs(ref[c][i] - conv.output[c][i])
			}
		}
	}

	core.Debug("conv max error:", max)
	if max > FFT_TOLERANCE {
		panic(fmt.Errorf("conv max error: %v > tolerance (%v)", max, FFT_TOLERANCE))
	}
	if max == 0 {
		// Something is wrong with the test
		panic(fmt.Errorf("conv max error: %v: too good to be true", max))
	}

}

// Run self-test and report estimated relative RMS error
// on forward+backward transform on data of magnitude order 1.
func (c *Conv1) Test() {
	N := c.n

	// make random input data between -1 and 1
	for i := 0; i < 3; i++ {
		for j := range c.input[i] {
			c.input[i][j] = 2*rand.Float32() - 1
		}
	}

	{
		c.noKernMul = true // don't do kernel multiplication, only fw+bw transform

		c.Push(N)
		c.Pull(N - 1)
		c.Pull(N)
		c.checkError()

		c.Push(N / 2)
		c.Push(N)
		c.Pull(N / 4)
		c.Pull(N - 1)
		c.Pull(N)
		c.checkError()

		c.noKernMul = false
	}

	// Set i/o arrays back to 0
	for i := 0; i < 3; i++ {
		for j := range c.input[i] {
			c.input[i][j] = 0
			c.output[i][j] = 0
		}
	}
}

const FFT_TOLERANCE = 1e-5 // panic if RMS error of fw+bw transform on random data is larger than this.

// Check if the rms error introduced by fw+bw transform is < FFT_TOLERANCE
func (c *Conv1) checkError() {
	NFFT := prod(PadSize(c.size))
	rms := 0.0
	for i := 0; i < 3; i++ {
		for j := range c.input[i] {
			rms += sqr(float64(c.input[i][j]) - float64(c.output[i][j])/float64(NFFT))
		}
	}
	rms = math.Sqrt(rms / float64(3*NFFT))

	core.Debug("RMS fft error:", rms)
	if rms > FFT_TOLERANCE {
		panic(fmt.Errorf("FFT RMS error: %v > tolerance (%v)", rms, FFT_TOLERANCE))
	}
	if rms == 0 {
		// Something is wrong with the test
		panic(fmt.Errorf("FFT RMS error: %v: too good to be true", rms))
	}
}

func sqr(x float64) float64 {
	return x * x
}