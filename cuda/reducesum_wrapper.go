package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var reducesum_code cu.Function

type reducesum_args struct {
	arg_src     unsafe.Pointer
	arg_dst     unsafe.Pointer
	arg_initVal float32
	arg_n       int
	argptr      [4]unsafe.Pointer
}

// Wrapper for reducesum CUDA kernel, asynchronous.
func k_reducesum_async(src unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int, cfg *Config, str cu.Stream) {
	if reducesum_code == 0 {
		reducesum_code = cu.ModuleLoadData(reducesum_ptx).GetFunction("reducesum")
	}

	var a reducesum_args

	a.arg_src = src
	a.argptr[0] = unsafe.Pointer(&a.arg_src)
	a.arg_dst = dst
	a.argptr[1] = unsafe.Pointer(&a.arg_dst)
	a.arg_initVal = initVal
	a.argptr[2] = unsafe.Pointer(&a.arg_initVal)
	a.arg_n = n
	a.argptr[3] = unsafe.Pointer(&a.arg_n)

	args := a.argptr[:]
	cu.LaunchKernel(reducesum_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, str, args)
}

// Wrapper for reducesum CUDA kernel, synchronized.
func k_reducesum(src unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int, cfg *Config) {
	str := Stream()
	k_reducesum_async(src, dst, initVal, n, cfg, str)
	SyncAndRecycle(str)
}

const reducesum_ptx = `
.version 3.0
.target sm_30
.address_size 64


.entry reducesum(
	.param .u64 reducesum_param_0,
	.param .u64 reducesum_param_1,
	.param .f32 reducesum_param_2,
	.param .u32 reducesum_param_3
)
{
	.reg .f32 	%f<32>;
	.reg .pred 	%p<8>;
	.reg .s32 	%r<48>;
	.reg .s64 	%rl<13>;
	// demoted variable
	.shared .align 4 .b8 __cuda_local_var_16163_32_non_const_sdata[2048];

	ld.param.u64 	%rl4, [reducesum_param_0];
	ld.param.u64 	%rl5, [reducesum_param_1];
	ld.param.u32 	%r1, [reducesum_param_3];
	cvta.to.global.u64 	%rl1, %rl5;
	cvta.to.global.u64 	%rl2, %rl4;
	.loc 2 8 1
	mov.u32 	%r2, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r3, %tid.x;
	mad.lo.s32 	%r46, %r2, %r11, %r3;
	mov.u32 	%r12, %nctaid.x;
	mul.lo.s32 	%r5, %r2, %r12;
	.loc 2 8 1
	setp.lt.s32 	%p1, %r46, %r1;
	@%p1 bra 	BB0_2;

	ld.param.f32 	%f31, [reducesum_param_2];
	bra.uni 	BB0_4;

BB0_2:
	ld.param.f32 	%f31, [reducesum_param_2];

BB0_3:
	.loc 2 8 1
	mul.wide.s32 	%rl6, %r46, 4;
	add.s64 	%rl7, %rl2, %rl6;
	ld.global.f32 	%f5, [%rl7];
	add.f32 	%f31, %f31, %f5;
	add.s32 	%r46, %r46, %r5;
	ld.param.u32 	%r37, [reducesum_param_3];
	.loc 2 8 1
	setp.lt.s32 	%p2, %r46, %r37;
	@%p2 bra 	BB0_3;

BB0_4:
	.loc 2 8 1
	mov.u32 	%r45, %tid.x;
	mul.wide.s32 	%rl8, %r45, 4;
	mov.u64 	%rl9, __cuda_local_var_16163_32_non_const_sdata;
	add.s64 	%rl3, %rl9, %rl8;
	.loc 2 8 1
	st.shared.f32 	[%rl3], %f31;
	bar.sync 	0;
	mov.u32 	%r42, %ntid.x;
	shr.u32 	%r47, %r42, 1;
	.loc 2 8 1
	setp.lt.u32 	%p3, %r47, 33;
	@%p3 bra 	BB0_8;

BB0_5:
	.loc 2 8 1
	setp.ge.u32 	%p4, %r3, %r47;
	@%p4 bra 	BB0_7;

	.loc 2 8 1
	ld.shared.f32 	%f6, [%rl3];
	add.s32 	%r16, %r47, %r3;
	mul.wide.u32 	%rl10, %r16, 4;
	add.s64 	%rl12, %rl9, %rl10;
	.loc 2 8 1
	ld.shared.f32 	%f7, [%rl12];
	add.f32 	%f8, %f6, %f7;
	st.shared.f32 	[%rl3], %f8;

BB0_7:
	.loc 2 8 1
	bar.sync 	0;
	shr.u32 	%r47, %r47, 1;
	.loc 2 8 1
	setp.gt.u32 	%p5, %r47, 32;
	@%p5 bra 	BB0_5;

BB0_8:
	.loc 2 8 1
	mov.u32 	%r44, %tid.x;
	.loc 2 8 1
	setp.gt.s32 	%p6, %r44, 31;
	@%p6 bra 	BB0_10;

	.loc 2 8 1
	ld.volatile.shared.f32 	%f9, [%rl3];
	ld.volatile.shared.f32 	%f10, [%rl3+128];
	add.f32 	%f11, %f9, %f10;
	st.volatile.shared.f32 	[%rl3], %f11;
	ld.volatile.shared.f32 	%f12, [%rl3+64];
	ld.volatile.shared.f32 	%f13, [%rl3];
	add.f32 	%f14, %f13, %f12;
	st.volatile.shared.f32 	[%rl3], %f14;
	ld.volatile.shared.f32 	%f15, [%rl3+32];
	ld.volatile.shared.f32 	%f16, [%rl3];
	add.f32 	%f17, %f16, %f15;
	st.volatile.shared.f32 	[%rl3], %f17;
	ld.volatile.shared.f32 	%f18, [%rl3+16];
	ld.volatile.shared.f32 	%f19, [%rl3];
	add.f32 	%f20, %f19, %f18;
	st.volatile.shared.f32 	[%rl3], %f20;
	ld.volatile.shared.f32 	%f21, [%rl3+8];
	ld.volatile.shared.f32 	%f22, [%rl3];
	add.f32 	%f23, %f22, %f21;
	st.volatile.shared.f32 	[%rl3], %f23;
	ld.volatile.shared.f32 	%f24, [%rl3+4];
	ld.volatile.shared.f32 	%f25, [%rl3];
	add.f32 	%f26, %f25, %f24;
	st.volatile.shared.f32 	[%rl3], %f26;

BB0_10:
	.loc 2 8 1
	mov.u32 	%r43, %tid.x;
	.loc 2 8 1
	setp.eq.s32 	%p7, %r43, 0;
	@%p7 bra 	BB0_12;

	.loc 2 9 2
	ret;

BB0_12:
	.loc 2 8 1
	ld.shared.f32 	%f27, [__cuda_local_var_16163_32_non_const_sdata];
	.loc 3 1844 5
	atom.global.add.f32 	%f28, [%rl1], %f27;
	.loc 2 9 2
	ret;
}


`