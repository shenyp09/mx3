package nc

type UniformSlice struct {
	value float32
	buf   []float32
}

func MakeUniformSlice() UniformSlice { return UniformSlice{0, []float32{0}} }

func (s *UniformSlice) SetValue(value float32) {
	s.value = value
}

func (s *UniformSlice) Range(i1, i2 int) []float32 {
	lenHave := len(s.buf)
	lenWant := i2 - i1
	if lenWant == 0 {
		return []float32{}
	}
	if lenHave == lenWant {
		if s.buf[0] != s.value {
			Memset(s.buf, s.value)
		}
		return s.buf
	}
	s.buf = ResizeBuffer(s.buf, lenWant)
	Memset(s.buf, s.value)
	return s.buf
}

func Memset(a []float32, value float32) {
	for i := range a {
		a[i] = value
	}
}

func ResizeBuffer(buf []float32, lenWant int) []float32 {
	lenHave := len(buf)
	switch {
	case lenHave == lenWant:
		return buf
	case lenHave > lenWant:
		return buf[:lenWant]
	case cap(buf) >= lenWant:
		newBuf := buf[:lenWant]
		return newBuf
	default:
		return make([]float32, lenWant)
	}
	return nil
}