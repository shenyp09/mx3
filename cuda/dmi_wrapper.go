package cuda

/*
 THIS FILE IS AUTO-GENERATED BY CUDA2GO.
 EDITING IS FUTILE.
*/

import (
	"github.com/barnex/cuda5/cu"
	"unsafe"
)

var dmi_code cu.Function

type dmi_args struct {
	arg_Hx unsafe.Pointer
	arg_Hy unsafe.Pointer
	arg_Hz unsafe.Pointer
	arg_mx unsafe.Pointer
	arg_my unsafe.Pointer
	arg_mz unsafe.Pointer
	arg_dx float32
	arg_dy float32
	arg_dz float32
	arg_cx float32
	arg_cy float32
	arg_cz float32
	arg_N0 int
	arg_N1 int
	arg_N2 int
	argptr [15]unsafe.Pointer
}

// Wrapper for dmi CUDA kernel, asynchronous.
func k_dmi_async(Hx unsafe.Pointer, Hy unsafe.Pointer, Hz unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, dx float32, dy float32, dz float32, cx float32, cy float32, cz float32, N0 int, N1 int, N2 int, cfg *Config, str cu.Stream) {
	if dmi_code == 0 {
		dmi_code = cu.ModuleLoadData(dmi_ptx).GetFunction("dmi")
	}

	var a dmi_args

	a.arg_Hx = Hx
	a.argptr[0] = unsafe.Pointer(&a.arg_Hx)
	a.arg_Hy = Hy
	a.argptr[1] = unsafe.Pointer(&a.arg_Hy)
	a.arg_Hz = Hz
	a.argptr[2] = unsafe.Pointer(&a.arg_Hz)
	a.arg_mx = mx
	a.argptr[3] = unsafe.Pointer(&a.arg_mx)
	a.arg_my = my
	a.argptr[4] = unsafe.Pointer(&a.arg_my)
	a.arg_mz = mz
	a.argptr[5] = unsafe.Pointer(&a.arg_mz)
	a.arg_dx = dx
	a.argptr[6] = unsafe.Pointer(&a.arg_dx)
	a.arg_dy = dy
	a.argptr[7] = unsafe.Pointer(&a.arg_dy)
	a.arg_dz = dz
	a.argptr[8] = unsafe.Pointer(&a.arg_dz)
	a.arg_cx = cx
	a.argptr[9] = unsafe.Pointer(&a.arg_cx)
	a.arg_cy = cy
	a.argptr[10] = unsafe.Pointer(&a.arg_cy)
	a.arg_cz = cz
	a.argptr[11] = unsafe.Pointer(&a.arg_cz)
	a.arg_N0 = N0
	a.argptr[12] = unsafe.Pointer(&a.arg_N0)
	a.arg_N1 = N1
	a.argptr[13] = unsafe.Pointer(&a.arg_N1)
	a.arg_N2 = N2
	a.argptr[14] = unsafe.Pointer(&a.arg_N2)

	args := a.argptr[:]
	cu.LaunchKernel(dmi_code, cfg.Grid.X, cfg.Grid.Y, cfg.Grid.Z, cfg.Block.X, cfg.Block.Y, cfg.Block.Z, 0, str, args)
}

// Wrapper for dmi CUDA kernel, synchronized.
func k_dmi(Hx unsafe.Pointer, Hy unsafe.Pointer, Hz unsafe.Pointer, mx unsafe.Pointer, my unsafe.Pointer, mz unsafe.Pointer, dx float32, dy float32, dz float32, cx float32, cy float32, cz float32, N0 int, N1 int, N2 int, cfg *Config) {
	str := Stream()
	k_dmi_async(Hx, Hy, Hz, mx, my, mz, dx, dy, dz, cx, cy, cz, N0, N1, N2, cfg, str)
	SyncAndRecycle(str)
}

const dmi_ptx = `
.version 3.0
.target sm_30
.address_size 64


.entry dmi(
	.param .u64 dmi_param_0,
	.param .u64 dmi_param_1,
	.param .u64 dmi_param_2,
	.param .u64 dmi_param_3,
	.param .u64 dmi_param_4,
	.param .u64 dmi_param_5,
	.param .f32 dmi_param_6,
	.param .f32 dmi_param_7,
	.param .f32 dmi_param_8,
	.param .f32 dmi_param_9,
	.param .f32 dmi_param_10,
	.param .f32 dmi_param_11,
	.param .u32 dmi_param_12,
	.param .u32 dmi_param_13,
	.param .u32 dmi_param_14
)
{
	.reg .f32 	%f<49>;
	.reg .pred 	%p<10>;
	.reg .s32 	%r<168>;
	.reg .s64 	%rl<53>;


	ld.param.u32 	%r1, [dmi_param_12];
	ld.param.u32 	%r2, [dmi_param_13];
	ld.param.u32 	%r3, [dmi_param_14];
	.loc 2 17 1
	mov.u32 	%r4, %ctaid.x;
	mov.u32 	%r12, %ntid.x;
	mov.u32 	%r5, %tid.x;
	mad.lo.s32 	%r6, %r12, %r4, %r5;
	.loc 2 18 1
	mov.u32 	%r13, %ntid.y;
	mov.u32 	%r14, %ctaid.y;
	mov.u32 	%r15, %tid.y;
	mad.lo.s32 	%r16, %r13, %r14, %r15;
	setp.lt.s32 	%p1, %r16, %r3;
	setp.lt.s32 	%p2, %r6, %r2;
	and.pred  	%p3, %p2, %p1;
	.loc 2 24 1
	setp.gt.s32 	%p4, %r1, 0;
	and.pred  	%p5, %p3, %p4;
	.loc 2 20 1
	@!%p5 bra 	BB0_9;

	ld.param.f32 	%f47, [dmi_param_10];
	.loc 2 32 1
	add.f32 	%f7, %f47, %f47;
	ld.param.f32 	%f48, [dmi_param_11];
	.loc 2 33 1
	add.f32 	%f8, %f48, %f48;
	ld.param.f32 	%f46, [dmi_param_9];
	.loc 2 38 1
	add.f32 	%f9, %f46, %f46;
	ld.param.u32 	%r163, [dmi_param_14];
	mad.lo.s32 	%r166, %r163, %r6, %r16;
	mov.u32 	%r17, 0;
	mov.u32 	%r167, %r17;

BB0_2:
	mov.u32 	%r9, %r167;
	ld.param.u64 	%rl44, [dmi_param_0];
	cvta.to.global.u64 	%rl10, %rl44;
	.loc 2 27 1
	mul.wide.s32 	%rl11, %r166, 4;
	add.s64 	%rl7, %rl10, %rl11;
	st.global.u32 	[%rl7], %r17;
	ld.param.u64 	%rl45, [dmi_param_1];
	cvta.to.global.u64 	%rl12, %rl45;
	.loc 2 28 1
	add.s64 	%rl8, %rl12, %rl11;
	st.global.u32 	[%rl8], %r17;
	ld.param.u64 	%rl46, [dmi_param_2];
	cvta.to.global.u64 	%rl13, %rl46;
	.loc 2 29 1
	add.s64 	%rl9, %rl13, %rl11;
	st.global.u32 	[%rl9], %r17;
	ld.param.f32 	%f41, [dmi_param_6];
	.loc 2 31 1
	setp.eq.f32 	%p6, %f41, 0f00000000;
	@%p6 bra 	BB0_4;

	mov.u32 	%r28, 0;
	.loc 3 238 5
	max.s32 	%r29, %r9, %r28;
	ld.param.u32 	%r154, [dmi_param_12];
	add.s32 	%r30, %r154, -1;
	.loc 3 210 5
	min.s32 	%r31, %r29, %r30;
	.loc 2 32 1
	add.s32 	%r36, %r6, 1;
	.loc 3 238 5
	max.s32 	%r37, %r36, %r28;
	ld.param.u32 	%r158, [dmi_param_13];
	add.s32 	%r38, %r158, -1;
	.loc 3 210 5
	min.s32 	%r39, %r37, %r38;
	mad.lo.s32 	%r40, %r31, %r158, %r39;
	.loc 3 238 5
	max.s32 	%r45, %r16, %r28;
	ld.param.u32 	%r162, [dmi_param_14];
	add.s32 	%r46, %r162, -1;
	.loc 3 210 5
	min.s32 	%r47, %r45, %r46;
	.loc 2 32 1
	mad.lo.s32 	%r48, %r40, %r162, %r47;
	ld.param.u64 	%rl52, [dmi_param_5];
	cvta.to.global.u64 	%rl14, %rl52;
	.loc 2 32 1
	mul.wide.s32 	%rl15, %r48, 4;
	add.s64 	%rl16, %rl14, %rl15;
	add.s32 	%r49, %r6, -1;
	.loc 3 238 5
	max.s32 	%r50, %r49, %r28;
	.loc 3 210 5
	min.s32 	%r51, %r50, %r38;
	mad.lo.s32 	%r52, %r31, %r158, %r51;
	add.s32 	%r53, %r16, -1;
	.loc 3 238 5
	max.s32 	%r54, %r53, %r28;
	.loc 3 210 5
	min.s32 	%r55, %r54, %r46;
	.loc 2 32 1
	mad.lo.s32 	%r56, %r52, %r162, %r55;
	mul.wide.s32 	%rl17, %r56, 4;
	add.s64 	%rl18, %rl14, %rl17;
	ld.global.f32 	%f10, [%rl18];
	ld.global.f32 	%f11, [%rl16];
	sub.f32 	%f12, %f11, %f10;
	.loc 4 1311 3
	div.rn.f32 	%f13, %f12, %f7;
	.loc 3 238 5
	max.s32 	%r59, %r6, %r28;
	.loc 3 210 5
	min.s32 	%r60, %r59, %r38;
	mad.lo.s32 	%r61, %r31, %r158, %r60;
	.loc 2 33 1
	add.s32 	%r62, %r16, 1;
	.loc 3 238 5
	max.s32 	%r63, %r62, %r28;
	.loc 3 210 5
	min.s32 	%r64, %r63, %r46;
	.loc 2 33 1
	mad.lo.s32 	%r65, %r61, %r162, %r64;
	ld.param.u64 	%rl50, [dmi_param_4];
	cvta.to.global.u64 	%rl19, %rl50;
	.loc 2 33 1
	mul.wide.s32 	%rl20, %r65, 4;
	add.s64 	%rl21, %rl19, %rl20;
	mad.lo.s32 	%r66, %r61, %r162, %r47;
	mul.wide.s32 	%rl22, %r66, 4;
	add.s64 	%rl23, %rl19, %rl22;
	ld.global.f32 	%f14, [%rl23];
	ld.global.f32 	%f15, [%rl21];
	sub.f32 	%f16, %f15, %f14;
	.loc 4 1311 3
	div.rn.f32 	%f17, %f16, %f8;
	sub.f32 	%f18, %f17, %f13;
	ld.param.f32 	%f40, [dmi_param_6];
	.loc 2 34 1
	mul.f32 	%f19, %f18, %f40;
	st.global.f32 	[%rl7], %f19;

BB0_4:
	ld.param.f32 	%f43, [dmi_param_7];
	.loc 2 37 1
	setp.eq.f32 	%p7, %f43, 0f00000000;
	@%p7 bra 	BB0_6;

	.loc 2 38 1
	add.s32 	%r70, %r9, 1;
	mov.u32 	%r71, 0;
	.loc 3 238 5
	max.s32 	%r72, %r70, %r71;
	ld.param.u32 	%r153, [dmi_param_12];
	add.s32 	%r73, %r153, -1;
	.loc 3 210 5
	min.s32 	%r74, %r72, %r73;
	ld.param.u32 	%r157, [dmi_param_13];
	add.s32 	%r75, %r157, -1;
	.loc 3 238 5
	max.s32 	%r76, %r6, %r71;
	.loc 3 210 5
	min.s32 	%r77, %r76, %r75;
	mad.lo.s32 	%r78, %r74, %r157, %r77;
	.loc 3 238 5
	max.s32 	%r83, %r16, %r71;
	ld.param.u32 	%r161, [dmi_param_14];
	add.s32 	%r84, %r161, -1;
	.loc 3 210 5
	min.s32 	%r85, %r83, %r84;
	.loc 2 38 1
	mad.lo.s32 	%r86, %r78, %r161, %r85;
	ld.param.u64 	%rl51, [dmi_param_5];
	cvta.to.global.u64 	%rl24, %rl51;
	.loc 2 38 1
	mul.wide.s32 	%rl25, %r86, 4;
	add.s64 	%rl26, %rl24, %rl25;
	add.s32 	%r87, %r9, -1;
	.loc 3 238 5
	max.s32 	%r88, %r87, %r71;
	.loc 3 210 5
	min.s32 	%r89, %r88, %r73;
	mad.lo.s32 	%r90, %r89, %r157, %r77;
	.loc 2 38 1
	mad.lo.s32 	%r91, %r90, %r161, %r85;
	mul.wide.s32 	%rl27, %r91, 4;
	add.s64 	%rl28, %rl24, %rl27;
	ld.global.f32 	%f20, [%rl28];
	ld.global.f32 	%f21, [%rl26];
	sub.f32 	%f22, %f21, %f20;
	.loc 4 1311 3
	div.rn.f32 	%f23, %f22, %f9;
	.loc 3 238 5
	max.s32 	%r94, %r9, %r71;
	.loc 3 210 5
	min.s32 	%r95, %r94, %r73;
	mad.lo.s32 	%r96, %r95, %r157, %r77;
	.loc 2 33 1
	add.s32 	%r97, %r16, 1;
	.loc 3 238 5
	max.s32 	%r98, %r97, %r71;
	.loc 3 210 5
	min.s32 	%r99, %r98, %r84;
	.loc 2 39 1
	mad.lo.s32 	%r100, %r96, %r161, %r99;
	ld.param.u64 	%rl48, [dmi_param_3];
	cvta.to.global.u64 	%rl29, %rl48;
	.loc 2 39 1
	mul.wide.s32 	%rl30, %r100, 4;
	add.s64 	%rl31, %rl29, %rl30;
	mad.lo.s32 	%r101, %r96, %r161, %r85;
	mul.wide.s32 	%rl32, %r101, 4;
	add.s64 	%rl33, %rl29, %rl32;
	ld.global.f32 	%f24, [%rl33];
	ld.global.f32 	%f25, [%rl31];
	sub.f32 	%f26, %f25, %f24;
	.loc 4 1311 3
	div.rn.f32 	%f27, %f26, %f8;
	.loc 2 40 1
	sub.f32 	%f28, %f23, %f27;
	ld.param.f32 	%f42, [dmi_param_7];
	.loc 2 40 1
	mul.f32 	%f29, %f28, %f42;
	st.global.f32 	[%rl8], %f29;

BB0_6:
	ld.param.f32 	%f45, [dmi_param_8];
	.loc 2 43 1
	setp.eq.f32 	%p8, %f45, 0f00000000;
	@%p8 bra 	BB0_8;

	.loc 2 44 1
	add.s32 	%r105, %r9, 1;
	mov.u32 	%r106, 0;
	.loc 3 238 5
	max.s32 	%r107, %r105, %r106;
	ld.param.u32 	%r152, [dmi_param_12];
	add.s32 	%r108, %r152, -1;
	.loc 3 210 5
	min.s32 	%r109, %r107, %r108;
	.loc 3 238 5
	max.s32 	%r114, %r6, %r106;
	ld.param.u32 	%r156, [dmi_param_13];
	add.s32 	%r115, %r156, -1;
	.loc 3 210 5
	min.s32 	%r116, %r114, %r115;
	mad.lo.s32 	%r117, %r109, %r156, %r116;
	.loc 3 238 5
	max.s32 	%r122, %r16, %r106;
	ld.param.u32 	%r160, [dmi_param_14];
	add.s32 	%r123, %r160, -1;
	.loc 3 210 5
	min.s32 	%r124, %r122, %r123;
	.loc 2 44 1
	mad.lo.s32 	%r125, %r117, %r160, %r124;
	ld.param.u64 	%rl49, [dmi_param_4];
	cvta.to.global.u64 	%rl34, %rl49;
	.loc 2 44 1
	mul.wide.s32 	%rl35, %r125, 4;
	add.s64 	%rl36, %rl34, %rl35;
	add.s32 	%r126, %r9, -1;
	.loc 3 238 5
	max.s32 	%r127, %r126, %r106;
	.loc 3 210 5
	min.s32 	%r128, %r127, %r108;
	mad.lo.s32 	%r129, %r128, %r156, %r116;
	.loc 2 44 1
	mad.lo.s32 	%r130, %r129, %r160, %r124;
	mul.wide.s32 	%rl37, %r130, 4;
	add.s64 	%rl38, %rl34, %rl37;
	ld.global.f32 	%f30, [%rl38];
	ld.global.f32 	%f31, [%rl36];
	sub.f32 	%f32, %f31, %f30;
	.loc 4 1311 3
	div.rn.f32 	%f33, %f32, %f9;
	.loc 3 238 5
	max.s32 	%r133, %r9, %r106;
	.loc 3 210 5
	min.s32 	%r134, %r133, %r108;
	.loc 2 32 1
	add.s32 	%r135, %r6, 1;
	.loc 3 238 5
	max.s32 	%r136, %r135, %r106;
	.loc 3 210 5
	min.s32 	%r137, %r136, %r115;
	mad.lo.s32 	%r138, %r134, %r156, %r137;
	.loc 2 45 1
	mad.lo.s32 	%r139, %r138, %r160, %r124;
	ld.param.u64 	%rl47, [dmi_param_3];
	cvta.to.global.u64 	%rl39, %rl47;
	.loc 2 45 1
	mul.wide.s32 	%rl40, %r139, 4;
	add.s64 	%rl41, %rl39, %rl40;
	add.s32 	%r140, %r6, -1;
	.loc 3 238 5
	max.s32 	%r141, %r140, %r106;
	.loc 3 210 5
	min.s32 	%r142, %r141, %r115;
	mad.lo.s32 	%r143, %r134, %r156, %r142;
	add.s32 	%r144, %r16, -1;
	.loc 3 238 5
	max.s32 	%r145, %r144, %r106;
	.loc 3 210 5
	min.s32 	%r146, %r145, %r123;
	.loc 2 45 1
	mad.lo.s32 	%r147, %r143, %r160, %r146;
	mul.wide.s32 	%rl42, %r147, 4;
	add.s64 	%rl43, %rl39, %rl42;
	ld.global.f32 	%f34, [%rl43];
	ld.global.f32 	%f35, [%rl41];
	sub.f32 	%f36, %f35, %f34;
	.loc 4 1311 3
	div.rn.f32 	%f37, %f36, %f7;
	sub.f32 	%f38, %f37, %f33;
	ld.param.f32 	%f44, [dmi_param_8];
	.loc 2 46 1
	mul.f32 	%f39, %f38, %f44;
	st.global.f32 	[%rl9], %f39;

BB0_8:
	.loc 2 24 18
	add.s32 	%r10, %r9, 1;
	ld.param.u32 	%r155, [dmi_param_13];
	ld.param.u32 	%r159, [dmi_param_14];
	mad.lo.s32 	%r166, %r159, %r155, %r166;
	ld.param.u32 	%r151, [dmi_param_12];
	.loc 2 24 1
	setp.lt.s32 	%p9, %r10, %r151;
	mov.u32 	%r167, %r10;
	.loc 2 24 1
	@%p9 bra 	BB0_2;

BB0_9:
	.loc 2 50 2
	ret;
}


`