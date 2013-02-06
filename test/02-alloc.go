// +build ignore

package main

import (
	"code.google.com/p/mx3/mx"
)

func main() {
	N0, N1, N2 := 2, 16, 32
	c0, c1, c2 := 1e-9, 1e-9, 1e-9
	mesh := mx.NewMesh(N0, N1, N2, c0, c1, c2)

	m := mx.NewQuant(mx.VECTOR, "m", "", mesh)
	m.Data().Memset(1, 0, 0)

	m.Data().WriteTo("m.mx3o")
}