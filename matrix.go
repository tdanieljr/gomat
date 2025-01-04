package gomat

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

func (m Matrix[T]) IterRows() func(yield func(int, []T) bool) {
	return func(yield func(int, []T) bool) {
		for i := range m.Rows {
			if !yield(i, m.Data[i*m.Stride:i*m.Stride+m.Cols]) {
				return
			}
		}
	}
}

func (m Matrix[T]) Sum() T {
	var tot T = 0
	for _, val := range m.Data {
		tot += val
	}
	return tot
}

func (m *Matrix[T]) Mul(other Matrix[T]) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	for i := range m.Data {
		m.Data[i] *= other.Data[i]
	}
}

func (m *Matrix[T]) Add(other Matrix[T]) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	for i := range m.Data {
		m.Data[i] += other.Data[i]
	}
}

func (m *Matrix[T]) Div(other Matrix[T]) {
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	for i := range m.Data {
		m.Data[i] /= other.Data[i]
	}
}

func (m *Matrix[T]) Apply(f func(T) T) {
	for i, val := range m.Data {
		m.Data[i] = f(val)
	}
}

func Linspace(start, stop, dx float64) []float64 {
	n := int64((stop-start)/dx) + 1
	x := make([]float64, n)
	for i := range n {
		x[i] = start + float64(i)*dx
	}
	return x
}

func RealToComplex(x []float64) []complex128 {
	y := make([]complex128, len(x))
	for i, val := range x {
		y[i] = complex(val, 0)
	}
	return y
}

func IntegrateRows[T Number](f Matrix[T], x []T) []T {
	// I don't think this is actualy effecient, we can likely just calculate this directly and save time
	integral := make([]T, f.Rows)
	for j, row := range f.IterRows() {
		integral[j] += Trapezoidal(row, x)
	}
	return integral
}

func Trapezoidal[T Number](f []T, x []T) T {
	var integral T
	for i := range len(x) - 1 {
		integral += (0.5) * (x[i+1] - x[i]) * (f[i+1] + f[i])
	}
	return integral
}
