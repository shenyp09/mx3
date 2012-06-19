package nc

import (
	"unsafe"
)

// Block is a [][][]float32 with square layout and contiguous underlying storage:
// 	len(block[0]) == len(block[1]) == ...
// 	len(block[i][0]) == len(block[i][1]) == ...
// 	len(block[i][j][0]) == len(block[i][j][1]) == ...
type Block struct {
	Array [][][]float32
	List  []float32
}

// Make a block of float32's of size N[0] x N[1] x N[2]
// with contiguous underlying storage.
func MakeBlock(size [3]int) Block {
	checkSize(size[:])
	storage := make([]float32, size[0]*size[1]*size[2])
	return Block{slice(storage, size), storage}
}

func slice(storage []float32, size [3]int) [][][]float32 {
	checkSize(size[:])
	sliced := make([][][]float32, size[0])
	for i := range sliced {
		sliced[i] = make([][]float32, size[1])
	}
	for i := range sliced {
		for j := range sliced[i] {
			sliced[i][j] = storage[(i*size[1]+j)*size[2]+0 : (i*size[1]+j)*size[2]+size[2]]
		}
	}
	return sliced
}

// Pointer to first element.
func (b *Block) Pointer() *float32 {
	return &b.List[0]
}

// unsafe.Pointer to first element.
func (b *Block) UnsafePointer() unsafe.Pointer {
	return unsafe.Pointer(b.Pointer())
}

// Total number of scalar elements.
func (b *Block) NFloat() int {
	a := b.Array
	return len(a) * len(a[0]) * len(a[0][0])
}

// Number of bytes used by underlying storage.
func (b *Block) Bytes() int64 {
	return SIZEOF_FLOAT32 * int64(b.NFloat())
}

// BlockSize is the size of the block (N0, N1, N2)
// as was passed to MakeBlock()
func (b *Block) BlockSize() [3]int {
	a := b.Array
	return [3]int{len(a), len(a[0]), len(a[0][0])}
}

// Is zero value?
func (b *Block) IsNil() bool {
	return b.List == nil
}

// Return the index'th slice of size WarpSize(),
// sharing underlying storage.
func (b *Block) Slice(index int) Block {
	storage := b.List[index*WarpLen() : (index+1)*WarpLen()]
	return Block{slice(storage, WarpSize()), storage}
}

const (
	SIZEOF_FLOAT32 = 4
)