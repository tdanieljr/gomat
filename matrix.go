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

// by Vec I mean a continuous 1d array of stride 1.
// The flat data behind a matrix or a slice of a matrix is not a vec.
// Each row should be when returned by IterRows
type Vec[T Number] []T

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
	m.indexCheck(i, j)
	return m.Data[i*m.Stride+j]
}

func (m *Matrix[T]) Set(i, j int, value T) {
	m.indexCheck(i, j)
	m.Data[i*m.Stride+j] = value
}

func (m Matrix[T]) Slice(i, k, j, l int) Matrix[T] {
	// Slice get the submatrix starting at i,j and going k rows and l columns from that point
	// The data in the slice points to the same backing data as the original matrix
	m.indexCheck(i, j)
	m.sliceCheck(i, k, j, l)
	return _slice(m, i, k, j, l, false)
}

func (m Matrix[T]) SliceWithCopy(i, k, j, l int) Matrix[T] {
	// same as Slice but the backing data is a new copy and disconnected from
	// the original matrix
	m.indexCheck(i, j)
	m.sliceCheck(i, k, j, l)
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

func (m Matrix[T]) IterRows() func(yield func(int, Vec[T]) bool) {
	return func(yield func(int, Vec[T]) bool) {
		for i := range m.Rows {
			if !yield(i, m.Data[i*m.Stride:i*m.Stride+m.Cols]) {
				return
			}
		}
	}
}

func (m Matrix[T]) Sum() T {
	var tot T = 0
	for _, row := range m.IterRows() {
		for _, val := range row {
			tot += val
		}
	}
	return tot
}

func (v Vec[T]) Sum() T {
	var tot T = 0
	for _, val := range v {
		tot += val
	}
	return tot
}

func (m Matrix[T]) GetRow(i int) Vec[T] {
	return m.Data[i*m.Stride : i*m.Stride+m.Cols]
}

func (m *Matrix[T]) Mul(other Matrix[T]) Matrix[T] {
	// TODO: what happens if strides are different? I think this is wrong.
	// since the backing data is flat
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	data := make([]T, m.Rows*m.Cols)
	for i, row := range m.IterRows() {
		other_row := other.GetRow(i)
		for j, val := range row {
			data[i*m.Cols+j] = val * other_row[j]
		}
	}
	return FromValues(m.Rows, m.Cols, data)
}

func (m *Matrix[T]) Div(other Matrix[T]) Matrix[T] {
	// TODO: what if strides are different?
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	data := make([]T, m.Rows*m.Cols)
	for i, row := range m.IterRows() {
		other_row := other.GetRow(i)
		for j, val := range row {
			data[i*m.Cols+j] = val + other_row[j]
		}
	}
	return FromValues(m.Rows, m.Cols, data)
}

func (m *Matrix[T]) Add(other Matrix[T]) Matrix[T] {
	// TODO: what if strides are different?
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	data := make([]T, m.Rows*m.Cols)
	for i, row := range m.IterRows() {
		other_row := other.GetRow(i)
		for j, val := range row {
			data[i*m.Cols+j] = val + other_row[j]
		}
	}
	return FromValues(m.Rows, m.Cols, data)
}

func (m *Matrix[T]) Sub(other Matrix[T]) Matrix[T] {
	// TODO: what if strides are different?
	if m.Rows != other.Rows || m.Cols != other.Cols {
		panic("mismatched dims")
	}
	data := make([]T, m.Rows*m.Cols)
	for i, row := range m.IterRows() {
		other_row := other.GetRow(i)
		for j, val := range row {
			data[i*m.Cols+j] = val - other_row[j]
		}
	}
	return FromValues(m.Rows, m.Cols, data)
}

func (m *Matrix[T]) Apply(f func(T) T) {
	for i, val := range m.Data {
		m.Data[i] = f(val)
	}
}

func Linspace(start, stop, dx float64) Vec[float64] {
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

func (m Matrix[T]) IntegrateRows(x Vec[T]) Vec[T] {
	// TODO: make sure this works with stride != Cols iterRows I think fixes it.
	integral := make(Vec[T], m.Rows)
	for j, row := range m.IterRows() {
		integral[j] += row.Trapezoidal(x)
	}
	return integral
}

func (v Vec[T]) Trapezoidal(x Vec[T]) T {
	var integral T
	for i := range len(x) - 1 {
		integral += (0.5) * (x[i+1] - x[i]) * (v[i+1] + v[i])
	}
	return integral
}

func (m Matrix[T]) indexCheck(i, j int) {
	if i >= m.Rows || j >= m.Cols {
		panic("Index out of bounds")
	}
	if i < 0 || j < 0 {
		panic("Error: Negative index")
	}
}

func (m Matrix[T]) sliceCheck(i, k, j, l int) {
	if (i+k) >= m.Rows || (j+l) >= m.Cols {
		panic("Slice goes out of bounds")
	}
	if k < 0 || l < 0 {
		panic("Error: Negative length slice")
	}
}
