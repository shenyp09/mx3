// 2D X (out-of-plane only) micromagnetic kernel multiplication:
// Mx = Kxx * Mx
//
// ~kernel has mirror symmetry along Y-axis,
// apart form first row,
// and is only stored (roughly) half:
//
// K00:
// xxxxx
// aaaaa
// bbbbb
// ....
// bbbbb
// aaaaa
//
extern "C" __global__ void 
kernmulRSymm2Dx(float* fftMx, float* fftKxx, int N1, int N2){

	int j = blockIdx.y * blockDim.y + threadIdx.y;
	int k = blockIdx.x * blockDim.x + threadIdx.x;

	if(j>= N1 || k>=N2){
 		return;	
	}

	int I = j*N2 + k;       // linear index for upper half of kernel
	int I2 = (N1-j)*N2 + k; // linear index for re-use of lower half

    float Kxx;

	if (j < N1/2 + 1){
		Kxx = fftKxx[I];
	}else{
		Kxx = fftKxx[I2];
	}

  	int e = 2 * I;

    float reMx = fftMx[e  ];
    float imMx = fftMx[e+1];

    fftMx[e  ] = reMx * Kxx;
    fftMx[e+1] = imMx * Kxx;
}
