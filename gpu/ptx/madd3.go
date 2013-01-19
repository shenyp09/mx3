package ptx

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"code.google.com/p/mx3/core"
	"github.com/barnex/cuda5/cu"
	"sync"
	"unsafe"
)

// pointers passed to CGO must be kept alive manually
// so we keep then here.
var (
	madd3_lock     sync.Mutex
	madd3_code     cu.Function
	madd3_stream   cu.Stream
	madd3_arg_dst  cu.DevicePtr
	madd3_arg_src1 cu.DevicePtr
	madd3_arg_fac1 float32
	madd3_arg_src2 cu.DevicePtr
	madd3_arg_fac2 float32
	madd3_arg_src3 cu.DevicePtr
	madd3_arg_fac3 float32
	madd3_arg_N    int

	madd3_argptr = [...]unsafe.Pointer{
		unsafe.Pointer(&madd3_arg_dst),
		unsafe.Pointer(&madd3_arg_src1),
		unsafe.Pointer(&madd3_arg_fac1),
		unsafe.Pointer(&madd3_arg_src2),
		unsafe.Pointer(&madd3_arg_fac2),
		unsafe.Pointer(&madd3_arg_src3),
		unsafe.Pointer(&madd3_arg_fac3),
		unsafe.Pointer(&madd3_arg_N)}
)

// CUDA kernel wrapper for madd3.
// The kernel is launched in a separate stream so that it can be parallel with memcpys etc.
// The stream is synchronized before this call returns.
func K_madd3(dst cu.DevicePtr, src1 cu.DevicePtr, fac1 float32, src2 cu.DevicePtr, fac2 float32, src3 cu.DevicePtr, fac3 float32, N int, gridDim, blockDim cu.Dim3) {
	madd3_lock.Lock()

	if madd3_stream == 0 {
		madd3_stream = cu.StreamCreate()
		core.Log("Loading PTX code for madd3")
		madd3_code = cu.ModuleLoadData(madd3_ptx).GetFunction("madd3")
	}

	madd3_arg_dst = dst
	madd3_arg_src1 = src1
	madd3_arg_fac1 = fac1
	madd3_arg_src2 = src2
	madd3_arg_fac2 = fac2
	madd3_arg_src3 = src3
	madd3_arg_fac3 = fac3
	madd3_arg_N = N

	args := madd3_argptr[:]
	cu.LaunchKernel(madd3_code, gridDim.X, gridDim.Y, gridDim.Z, blockDim.X, blockDim.Y, blockDim.Z, 0, madd3_stream, args)
	madd3_stream.Synchronize()
	madd3_lock.Unlock()
}

const madd3_ptx = `
.version 3.1
.target sm_30
.address_size 64


.visible .entry madd3(
	.param .u64 madd3_param_0,
	.param .u64 madd3_param_1,
	.param .f32 madd3_param_2,
	.param .u64 madd3_param_3,
	.param .f32 madd3_param_4,
	.param .u64 madd3_param_5,
	.param .f32 madd3_param_6,
	.param .u32 madd3_param_7
)
{
	.reg .pred 	%p<2>;
	.reg .s32 	%r<13>;
	.reg .f32 	%f<10>;
	.reg .s64 	%rd<14>;


	ld.param.u64 	%rd5, [madd3_param_0];
	ld.param.u64 	%rd6, [madd3_param_1];
	ld.param.f32 	%f1, [madd3_param_2];
	ld.param.u64 	%rd7, [madd3_param_3];
	ld.param.f32 	%f2, [madd3_param_4];
	ld.param.u64 	%rd8, [madd3_param_5];
	ld.param.f32 	%f3, [madd3_param_6];
	ld.param.u32 	%r2, [madd3_param_7];
	cvta.to.global.u64 	%rd1, %rd5;
	cvta.to.global.u64 	%rd2, %rd8;
	cvta.to.global.u64 	%rd3, %rd7;
	cvta.to.global.u64 	%rd4, %rd6;
	.loc 2 4 1
	mov.u32 	%r3, %nctaid.x;
	mov.u32 	%r4, %ctaid.y;
	mov.u32 	%r5, %ctaid.x;
	mad.lo.s32 	%r6, %r3, %r4, %r5;
	mov.u32 	%r7, %ntid.x;
	mov.u32 	%r8, %tid.x;
	mad.lo.s32 	%r1, %r6, %r7, %r8;
	.loc 2 5 1
	setp.ge.s32 	%p1, %r1, %r2;
	@%p1 bra 	BB0_2;

	.loc 2 6 1
	mul.wide.s32 	%rd9, %r1, 4;
	add.s64 	%rd10, %rd4, %rd9;
	ld.global.f32 	%f4, [%rd10];
	add.s64 	%rd11, %rd3, %rd9;
	ld.global.f32 	%f5, [%rd11];
	add.s64 	%rd12, %rd2, %rd9;
	ld.global.f32 	%f6, [%rd12];
	mul.f32 	%f7, %f6, %f3;
	fma.rn.f32 	%f8, %f5, %f2, %f7;
	fma.rn.f32 	%f9, %f4, %f1, %f8;
	add.s64 	%rd13, %rd1, %rd9;
	st.global.f32 	[%rd13], %f9;

BB0_2:
	.loc 2 9 2
	ret;
}


`