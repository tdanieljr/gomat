package mat

type Number interface {
	float64 | complex128
}

type Matrix[T Number] struct {
	Data   []T
	Rows   int
	Cols   int
	Stride int
}

func Init[T Number](rows, cols int) Matrix[T] {
	data := make([]T, rows*cols)
	return Matrix[T]{Data: data, Rows: rows, Cols: cols, Stride: cols}
}

func FromValues[T Number](rows, cols int, data []T) Matrix[T] {
	if len(data) != rows*cols {
		panic("mismatch data length and rowxcol count")
	}
	return Matrix[T]{Data: data, Rows: rows, Cols: cols, Stride: cols}
}

func (m Matrix[T]) Get(i, j int) T {
	return m.Data[i*m.Stride+j]
}

func (m *Matrix[T]) Set(i, j int, value T) {
	m.Data[i*m.Stride+j] = value
}

func (m Matrix[T]) Slice(i, k, j, l int) Matrix[T] {
	return _slice(m, i, k, j, l, false)
}

func (m Matrix[T]) SliceWithCopy(i, k, j, l int) Matrix[T] {
	return _slice(m, i, k, j, l, true)
}

func _slice[T Number](m Matrix[T], i, k, j, l int, clone bool) Matrix[T] {
	start := i*m.Stride + j
	n := (k-1)*m.Stride + l

	if clone {
		data := make([]T, n)
		copy(data, m.Data[start:start+n])
		return Matrix[T]{data, k, l, m.Stride}
	}
	return Matrix[T]{m.Data[start : start+n], k, l, m.Stride}
}
