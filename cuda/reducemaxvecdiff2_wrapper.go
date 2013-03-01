package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var reducemaxvecdiff2_code cu.Function

type reducemaxvecdiff2_args struct {
	arg_x1      unsafe.Pointer
	arg_y1      unsafe.Pointer
	arg_z1      unsafe.Pointer
	arg_x2      unsafe.Pointer
	arg_y2      unsafe.Pointer
	arg_z2      unsafe.Pointer
	arg_dst     unsafe.Pointer
	arg_initVal float32
	arg_n       int
	argptr      [9]unsafe.Pointer
}

// Wrapper for reducemaxvecdiff2 CUDA kernel, asynchronous.
func k_reducemaxvecdiff2_async(x1 unsafe.Pointer, y1 unsafe.Pointer, z1 unsafe.Pointer, x2 unsafe.Pointer, y2 unsafe.Pointer, z2 unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int, cfg *Config, str cu.Stream) {
	if reducemaxvecdiff2_code == 0 {
		reducemaxvecdiff2_code = cu.ModuleLoadData(reducemaxvecdiff2_ptx).GetFunction("reducemaxvecdiff2")
	}

	var a reducemaxvecdiff2_args

	a.arg_x1 = x1
	a.argptr[0] = unsafe.Pointer(&a.arg_x1)
	a.arg_y1 = y1
	a.argptr[1] = unsafe.Pointer(&a.arg_y1)
	a.arg_z1 = z1
	a.argptr[2] = unsafe.Pointer(&a.arg_z1)
	a.arg_x2 = x2
	a.argptr[3] = unsafe.Pointer(&a.arg_x2)
	a.arg_y2 = y2
	a.argptr[4] = unsafe.Pointer(&a.arg_y2)
	a.arg_z2 = z2
	a.argptr[5] = unsafe.Pointer(&a.arg_z2)
	a.arg_dst = dst
	a.argptr[6] = unsafe.Pointer(&a.arg_dst)
	a.arg_initVal = initVal
	a.argptr[7] = unsafe.Pointer(&a.arg_initVal)
	a.arg_n = n
	a.argptr[8] = unsafe.Pointer(&a.arg_n)

	args := a.argptr[:]
	cu.LaunchKernel(reducemaxvecdiff2_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, str, args)
}

// Wrapper for reducemaxvecdiff2 CUDA kernel, synchronized.
func k_reducemaxvecdiff2(x1 unsafe.Pointer, y1 unsafe.Pointer, z1 unsafe.Pointer, x2 unsafe.Pointer, y2 unsafe.Pointer, z2 unsafe.Pointer, dst unsafe.Pointer, initVal float32, n int, cfg *Config) {
	str := Stream()
	k_reducemaxvecdiff2_async(x1, y1, z1, x2, y2, z2, dst, initVal, n, cfg, str)
	SyncAndRecycle(str)
}

const reducemaxvecdiff2_ptx = `
.version 3.0
.target sm_30
.address_size 64


.entry reducemaxvecdiff2(
	.param .u64 reducemaxvecdiff2_param_0,
	.param .u64 reducemaxvecdiff2_param_1,
	.param .u64 reducemaxvecdiff2_param_2,
	.param .u64 reducemaxvecdiff2_param_3,
	.param .u64 reducemaxvecdiff2_param_4,
	.param .u64 reducemaxvecdiff2_param_5,
	.param .u64 reducemaxvecdiff2_param_6,
	.param .f32 reducemaxvecdiff2_param_7,
	.param .u32 reducemaxvecdiff2_param_8
)
{
	.reg .f32 	%f<43>;
	.reg .pred 	%p<8>;
	.reg .s32 	%r<55>;
	.reg .s64 	%rl<29>;
	// demoted variable
	.shared .align 4 .b8 __cuda_local_var_16238_32_non_const_sdata[2048];

	ld.param.u64 	%rl9, [reducemaxvecdiff2_param_0];
	ld.param.u64 	%rl10, [reducemaxvecdiff2_param_1];
	ld.param.u64 	%rl11, [reducemaxvecdiff2_param_2];
	ld.param.u64 	%rl12, [reducemaxvecdiff2_param_3];
	ld.param.u64 	%rl13, [reducemaxvecdiff2_param_4];
	ld.param.u64 	%rl14, [reducemaxvecdiff2_param_5];
	ld.param.u32 	%r1, [reducemaxvecdiff2_param_8];
	cvta.to.global.u64 	%rl2, %rl14;
	cvta.to.global.u64 	%rl3, %rl11;
	cvta.to.global.u64 	%rl4, %rl13;
	cvta.to.global.u64 	%rl5, %rl10;
	cvta.to.global.u64 	%rl6, %rl12;
	cvta.to.global.u64 	%rl7, %rl9;
	.loc 2 14 1
	mov.u32 	%r2, %ntid.x;
	mov.u32 	%r11, %ctaid.x;
	mov.u32 	%r3, %tid.x;
	mad.lo.s32 	%r53, %r2, %r11, %r3;
	mov.u32 	%r12, %nctaid.x;
	mul.lo.s32 	%r5, %r2, %r12;
	.loc 2 14 1
	setp.lt.s32 	%p1, %r53, %r1;
	@%p1 bra 	BB0_2;

	ld.param.f32 	%f42, [reducemaxvecdiff2_param_7];
	bra.uni 	BB0_4;

BB0_2:
	ld.param.f32 	%f42, [reducemaxvecdiff2_param_7];

BB0_3:
	.loc 2 14 1
	mul.wide.s32 	%rl15, %r53, 4;
	add.s64 	%rl16, %rl7, %rl15;
	add.s64 	%rl17, %rl6, %rl15;
	ld.global.f32 	%f5, [%rl17];
	ld.global.f32 	%f6, [%rl16];
	sub.f32 	%f7, %f6, %f5;
	add.s64 	%rl18, %rl5, %rl15;
	add.s64 	%rl19, %rl4, %rl15;
	ld.global.f32 	%f8, [%rl19];
	ld.global.f32 	%f9, [%rl18];
	sub.f32 	%f10, %f9, %f8;
	mul.f32 	%f11, %f10, %f10;
	fma.rn.f32 	%f12, %f7, %f7, %f11;
	add.s64 	%rl20, %rl3, %rl15;
	add.s64 	%rl21, %rl2, %rl15;
	ld.global.f32 	%f13, [%rl21];
	ld.global.f32 	%f14, [%rl20];
	sub.f32 	%f15, %f14, %f13;
	fma.rn.f32 	%f16, %f15, %f15, %f12;
	.loc 3 435 5
	max.f32 	%f42, %f42, %f16;
	.loc 2 14 1
	add.s32 	%r53, %r53, %r5;
	ld.param.u32 	%r44, [reducemaxvecdiff2_param_8];
	.loc 2 14 1
	setp.lt.s32 	%p2, %r53, %r44;
	@%p2 bra 	BB0_3;

BB0_4:
	.loc 2 14 1
	mov.u32 	%r52, %tid.x;
	mul.wide.s32 	%rl22, %r52, 4;
	mov.u64 	%rl23, __cuda_local_var_16238_32_non_const_sdata;
	add.s64 	%rl8, %rl23, %rl22;
	.loc 2 14 1
	st.shared.f32 	[%rl8], %f42;
	bar.sync 	0;
	mov.u32 	%r49, %ntid.x;
	shr.u32 	%r54, %r49, 1;
	.loc 2 14 1
	setp.lt.u32 	%p3, %r54, 33;
	@%p3 bra 	BB0_8;

BB0_5:
	.loc 2 14 1
	setp.ge.u32 	%p4, %r3, %r54;
	@%p4 bra 	BB0_7;

	.loc 2 14 1
	ld.shared.f32 	%f17, [%rl8];
	add.s32 	%r21, %r54, %r3;
	mul.wide.u32 	%rl24, %r21, 4;
	add.s64 	%rl26, %rl23, %rl24;
	.loc 2 14 1
	ld.shared.f32 	%f18, [%rl26];
	.loc 3 435 5
	max.f32 	%f19, %f17, %f18;
	.loc 2 14 1
	st.shared.f32 	[%rl8], %f19;

BB0_7:
	.loc 2 14 1
	bar.sync 	0;
	shr.u32 	%r54, %r54, 1;
	.loc 2 14 1
	setp.gt.u32 	%p5, %r54, 32;
	@%p5 bra 	BB0_5;

BB0_8:
	.loc 2 14 1
	mov.u32 	%r51, %tid.x;
	.loc 2 14 1
	setp.gt.s32 	%p6, %r51, 31;
	@%p6 bra 	BB0_10;

	.loc 2 14 1
	ld.volatile.shared.f32 	%f20, [%rl8];
	ld.volatile.shared.f32 	%f21, [%rl8+128];
	.loc 3 435 5
	max.f32 	%f22, %f20, %f21;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f22;
	ld.volatile.shared.f32 	%f23, [%rl8+64];
	ld.volatile.shared.f32 	%f24, [%rl8];
	.loc 3 435 5
	max.f32 	%f25, %f24, %f23;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f25;
	ld.volatile.shared.f32 	%f26, [%rl8+32];
	ld.volatile.shared.f32 	%f27, [%rl8];
	.loc 3 435 5
	max.f32 	%f28, %f27, %f26;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f28;
	ld.volatile.shared.f32 	%f29, [%rl8+16];
	ld.volatile.shared.f32 	%f30, [%rl8];
	.loc 3 435 5
	max.f32 	%f31, %f30, %f29;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f31;
	ld.volatile.shared.f32 	%f32, [%rl8+8];
	ld.volatile.shared.f32 	%f33, [%rl8];
	.loc 3 435 5
	max.f32 	%f34, %f33, %f32;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f34;
	ld.volatile.shared.f32 	%f35, [%rl8+4];
	ld.volatile.shared.f32 	%f36, [%rl8];
	.loc 3 435 5
	max.f32 	%f37, %f36, %f35;
	.loc 2 14 1
	st.volatile.shared.f32 	[%rl8], %f37;

BB0_10:
	.loc 2 14 1
	mov.u32 	%r50, %tid.x;
	.loc 2 14 1
	setp.eq.s32 	%p7, %r50, 0;
	@%p7 bra 	BB0_12;

	.loc 2 15 2
	ret;

BB0_12:
	.loc 2 14 1
	ld.shared.f32 	%f38, [__cuda_local_var_16238_32_non_const_sdata];
	.loc 3 395 5
	abs.f32 	%f39, %f38;
	mov.b32 	 %r42, %f39;
	ld.param.u64 	%rl28, [reducemaxvecdiff2_param_6];
	cvta.to.global.u64 	%rl27, %rl28;
	.loc 3 1881 5
	atom.global.max.s32 	%r43, [%rl27], %r42;
	.loc 2 15 2
	ret;
}


`