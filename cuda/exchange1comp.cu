#include "stencil.h"

// Add 1 component of exchange interaction to Beff (Tesla).
// m is normalized.

extern "C" __global__ void
addexchange1comp(float* __restrict__ Beff, float* __restrict__ m, 
                 float wx, float wy, float wz, int N0, int N1, int N2){

	int j = blockIdx.x * blockDim.x + threadIdx.x;
	int k = blockIdx.y * blockDim.y + threadIdx.y;

	if (j >= N1 || k >= N2){
		return;
	}

	for(int i=0; i<N0; i++){

		int I = idx(i, j, k);
		float B = Beff[I];
		float m0 = m[I];

		float m1 = m[idx(i, j, lclamp(k-1    ))];
		float m2 = m[idx(i, j, hclamp(k+1, N2))];
		B += wz * ((m1-m0) + (m2-m0));

		m1 = m[idx(i, lclamp(j-1   ), k)];
		m2 = m[idx(i, hclamp(j+1,N1), k)];
		B += wy * ((m1-m0) + (m2-m0));

		// only take vertical derivative for 3D sim
		if (N0 != 1){
			m1 = m[idx(hclamp(i+1,N0), j, k)];
			m2 = m[idx(lclamp(i-1   ), j, k)];
			B  += wx * ((m1-m0) + (m2-m0));
		}

		Beff[I] = B;
	}
}

