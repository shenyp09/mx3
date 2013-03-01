package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var copypad_code cu.Function

type copypad_args struct {
	arg_dst unsafe.Pointer
	arg_D0  int
	arg_D1  int
	arg_D2  int
	arg_src unsafe.Pointer
	arg_S0  int
	arg_S1  int
	arg_S2  int
	argptr  [8]unsafe.Pointer
}

// Wrapper for copypad CUDA kernel, asynchronous.
func k_copypad_async(dst unsafe.Pointer, D0 int, D1 int, D2 int, src unsafe.Pointer, S0 int, S1 int, S2 int, cfg *Config, str cu.Stream) {
	if copypad_code == 0 {
		copypad_code = cu.ModuleLoadData(copypad_ptx).GetFunction("copypad")
	}

	var a copypad_args

	a.arg_dst = dst
	a.argptr[0] = unsafe.Pointer(&a.arg_dst)
	a.arg_D0 = D0
	a.argptr[1] = unsafe.Pointer(&a.arg_D0)
	a.arg_D1 = D1
	a.argptr[2] = unsafe.Pointer(&a.arg_D1)
	a.arg_D2 = D2
	a.argptr[3] = unsafe.Pointer(&a.arg_D2)
	a.arg_src = src
	a.argptr[4] = unsafe.Pointer(&a.arg_src)
	a.arg_S0 = S0
	a.argptr[5] = unsafe.Pointer(&a.arg_S0)
	a.arg_S1 = S1
	a.argptr[6] = unsafe.Pointer(&a.arg_S1)
	a.arg_S2 = S2
	a.argptr[7] = unsafe.Pointer(&a.arg_S2)

	args := a.argptr[:]
	cu.LaunchKernel(copypad_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, str, args)
}

// Wrapper for copypad CUDA kernel, synchronized.
func k_copypad(dst unsafe.Pointer, D0 int, D1 int, D2 int, src unsafe.Pointer, S0 int, S1 int, S2 int, cfg *Config) {
	str := Stream()
	k_copypad_async(dst, D0, D1, D2, src, S0, S1, S2, cfg, str)
	SyncAndRecycle(str)
}

const copypad_ptx = `
.version 3.0
.target sm_30
.address_size 64


.entry copypad(
	.param .u64 copypad_param_0,
	.param .u32 copypad_param_1,
	.param .u32 copypad_param_2,
	.param .u32 copypad_param_3,
	.param .u64 copypad_param_4,
	.param .u32 copypad_param_5,
	.param .u32 copypad_param_6,
	.param .u32 copypad_param_7
)
{
	.reg .f32 	%f<2>;
	.reg .pred 	%p<10>;
	.reg .s32 	%r<36>;
	.reg .s64 	%rl<25>;


	ld.param.u64 	%rl11, [copypad_param_0];
	ld.param.u32 	%r2, [copypad_param_2];
	ld.param.u32 	%r3, [copypad_param_3];
	ld.param.u64 	%rl12, [copypad_param_4];
	ld.param.u32 	%r5, [copypad_param_6];
	ld.param.u32 	%r6, [copypad_param_7];
	cvta.to.global.u64 	%rl1, %rl11;
	cvta.to.global.u64 	%rl2, %rl12;
	.loc 2 11 1
	mov.u32 	%r7, %ntid.y;
	mov.u32 	%r8, %ctaid.y;
	mov.u32 	%r9, %tid.y;
	mad.lo.s32 	%r16, %r7, %r8, %r9;
	.loc 2 12 1
	mov.u32 	%r10, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r12, %tid.x;
	mad.lo.s32 	%r17, %r10, %r11, %r12;
	.loc 2 14 1
	setp.ge.s32 	%p1, %r17, %r6;
	setp.ge.s32 	%p2, %r16, %r5;
	or.pred  	%p3, %p1, %p2;
	.loc 2 14 1
	setp.ge.s32 	%p4, %r16, %r2;
	or.pred  	%p5, %p3, %p4;
	.loc 2 14 1
	setp.ge.s32 	%p6, %r17, %r3;
	or.pred  	%p7, %p5, %p6;
	.loc 2 14 1
	@%p7 bra 	BB0_4;

	ld.param.u32 	%r23, [copypad_param_1];
	ld.param.u32 	%r26, [copypad_param_5];
	.loc 3 210 5
	min.s32 	%r13, %r26, %r23;
	.loc 2 20 1
	setp.lt.s32 	%p8, %r13, 1;
	@%p8 bra 	BB0_4;

	cvt.s64.s32 	%rl13, %r16;
	ld.param.u32 	%r25, [copypad_param_3];
	cvt.s64.s32 	%rl14, %r25;
	cvt.s64.s32 	%rl15, %r17;
	mad.lo.s64 	%rl16, %rl14, %rl13, %rl15;
	shl.b64 	%rl17, %rl16, 2;
	add.s64 	%rl24, %rl1, %rl17;
	ld.param.u32 	%r24, [copypad_param_2];
	mul.wide.s32 	%rl18, %r25, %r24;
	shl.b64 	%rl4, %rl18, 2;
	ld.param.u32 	%r28, [copypad_param_7];
	cvt.s64.s32 	%rl19, %r28;
	mad.lo.s64 	%rl20, %rl19, %rl13, %rl15;
	shl.b64 	%rl21, %rl20, 2;
	add.s64 	%rl23, %rl2, %rl21;
	ld.param.u32 	%r27, [copypad_param_6];
	mul.wide.s32 	%rl22, %r28, %r27;
	shl.b64 	%rl6, %rl22, 2;
	mov.u32 	%r35, 0;

BB0_3:
	.loc 2 21 1
	ld.global.f32 	%f1, [%rl23];
	st.global.f32 	[%rl24], %f1;
	add.s64 	%rl24, %rl24, %rl4;
	add.s64 	%rl23, %rl23, %rl6;
	.loc 2 20 52
	add.s32 	%r35, %r35, 1;
	.loc 2 20 1
	setp.lt.s32 	%p9, %r35, %r13;
	@%p9 bra 	BB0_3;

BB0_4:
	.loc 2 23 2
	ret;
}


`