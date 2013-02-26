package kernel

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var normalize_code cu.Function

type normalize_args struct {
	arg_vx   unsafe.Pointer
	arg_vy   unsafe.Pointer
	arg_vz   unsafe.Pointer
	arg_mask unsafe.Pointer
	arg_norm float32
	arg_N    int
	argptr   [6]unsafe.Pointer
}

// Wrapper for normalize CUDA kernel, asynchronous.
func K_normalize_async(vx unsafe.Pointer, vy unsafe.Pointer, vz unsafe.Pointer, mask unsafe.Pointer, norm float32, N int, gridDim, blockDim cu.Dim3, str cu.Stream) {
	if normalize_code == 0 {
		normalize_code = cu.ModuleLoadData(normalize_ptx).GetFunction("normalize")
	}

	var a normalize_args

	a.arg_vx = vx
	a.argptr[0] = unsafe.Pointer(&a.arg_vx)
	a.arg_vy = vy
	a.argptr[1] = unsafe.Pointer(&a.arg_vy)
	a.arg_vz = vz
	a.argptr[2] = unsafe.Pointer(&a.arg_vz)
	a.arg_mask = mask
	a.argptr[3] = unsafe.Pointer(&a.arg_mask)
	a.arg_norm = norm
	a.argptr[4] = unsafe.Pointer(&a.arg_norm)
	a.arg_N = N
	a.argptr[5] = unsafe.Pointer(&a.arg_N)

	args := a.argptr[:]
	cu.LaunchKernel(normalize_code, gridDim.X, gridDim.Y, gridDim.Z, blockDim.X, blockDim.Y, blockDim.Z, 0, str, args)
}

// Wrapper for normalize CUDA kernel, synchronized.
func K_normalize(vx unsafe.Pointer, vy unsafe.Pointer, vz unsafe.Pointer, mask unsafe.Pointer, norm float32, N int, gridDim, blockDim cu.Dim3) {
	str := Stream()
	K_normalize_async(vx, vy, vz, mask, norm, N, gridDim, blockDim, str)
	SyncAndRecycle(str)
}

const normalize_ptx = `
.version 3.1
.target sm_30
.address_size 64


.visible .entry normalize(
	.param .u64 normalize_param_0,
	.param .u64 normalize_param_1,
	.param .u64 normalize_param_2,
	.param .u64 normalize_param_3,
	.param .f32 normalize_param_4,
	.param .u32 normalize_param_5
)
{
	.reg .pred 	%p<4>;
	.reg .s32 	%r<16>;
	.reg .f32 	%f<23>;
	.reg .s64 	%rd<16>;


	ld.param.u64 	%rd10, [normalize_param_0];
	ld.param.u64 	%rd11, [normalize_param_1];
	ld.param.u64 	%rd12, [normalize_param_2];
	ld.param.u64 	%rd9, [normalize_param_3];
	ld.param.f32 	%f21, [normalize_param_4];
	ld.param.u32 	%r2, [normalize_param_5];
	cvta.to.global.u64 	%rd1, %rd9;
	cvta.to.global.u64 	%rd2, %rd12;
	cvta.to.global.u64 	%rd3, %rd11;
	cvta.to.global.u64 	%rd4, %rd10;
	.loc 2 7 1
	mov.u32 	%r3, %nctaid.x;
	mov.u32 	%r4, %ctaid.y;
	mov.u32 	%r5, %ctaid.x;
	mad.lo.s32 	%r6, %r3, %r4, %r5;
	mov.u32 	%r7, %ntid.x;
	mov.u32 	%r8, %tid.x;
	mad.lo.s32 	%r1, %r6, %r7, %r8;
	.loc 2 8 1
	setp.ge.s32 	%p1, %r1, %r2;
	@%p1 bra 	BB0_6;

	.loc 2 10 1
	cvt.s64.s32 	%rd5, %r1;
	mul.wide.s32 	%rd13, %r1, 4;
	add.s64 	%rd6, %rd4, %rd13;
	ld.global.f32 	%f1, [%rd6];
	add.s64 	%rd7, %rd3, %rd13;
	ld.global.f32 	%f2, [%rd7];
	add.s64 	%rd8, %rd2, %rd13;
	ld.global.f32 	%f3, [%rd8];
	.loc 2 11 1
	setp.eq.s64 	%p2, %rd9, 0;
	@%p2 bra 	BB0_3;

	shl.b64 	%rd14, %rd5, 2;
	add.s64 	%rd15, %rd1, %rd14;
	ld.global.f32 	%f10, [%rd15];
	mul.f32 	%f21, %f10, %f21;

BB0_3:
	.loc 2 12 1
	mul.f32 	%f12, %f2, %f2;
	fma.rn.f32 	%f13, %f1, %f1, %f12;
	fma.rn.f32 	%f14, %f3, %f3, %f13;
	.loc 3 991 5
	sqrt.rn.f32 	%f6, %f14;
	mov.f32 	%f22, 0f00000000;
	.loc 2 12 1
	setp.eq.f32 	%p3, %f6, 0f00000000;
	@%p3 bra 	BB0_5;

	rcp.rn.f32 	%f22, %f6;

BB0_5:
	mul.f32 	%f15, %f22, %f1;
	mul.f32 	%f16, %f21, %f15;
	mul.f32 	%f17, %f22, %f2;
	mul.f32 	%f18, %f21, %f17;
	mul.f32 	%f19, %f22, %f3;
	mul.f32 	%f20, %f21, %f19;
	.loc 2 13 1
	st.global.f32 	[%rd6], %f16;
	.loc 2 14 1
	st.global.f32 	[%rd7], %f18;
	.loc 2 15 1
	st.global.f32 	[%rd8], %f20;

BB0_6:
	.loc 2 17 2
	ret;
}


`