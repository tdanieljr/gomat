package gomat

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMatrix(t *testing.T) {
	t.Run("testing init", func(t *testing.T) {
		got := Init[float64](2, 2)
		want := Matrix[float64]{[]float64{0, 0, 0, 0}, 2, 2, 2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("init not as expected")
		}
	})
	t.Run("testing from values", func(t *testing.T) {
		got := FromValues(2, 2, []float64{0, 0, 0, 0})
		want := Matrix[float64]{[]float64{0, 0, 0, 0}, 2, 2, 2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("init not as expected")
		}
	})
	t.Run("testing get/set", func(t *testing.T) {
		m := Init[float64](20, 20)
		if m.Get(3, 7) != 0 {
			t.Errorf("Get wrong")
		}
		m.Set(3, 7, 100)
		if m.Get(3, 7) != 100 {
			t.Errorf("Set wrong")
		}
	})
	t.Run("testing slice", func(t *testing.T) {
		m := Init[float64](10, 10)
		m.Set(0, 5, 10)
		a := m.Slice(0, 3, 2, 7)
		if a.Get(0, 3) != 10 {
			t.Errorf("Slice wrong")
		}
	})
	t.Run("testing iter", func(t *testing.T) {
		// since the backing data is flat I slice the matrix
		// to test if the iter method return contiguous rows even if the
		// stride is not equal to the cols
		want := Vec[float64]{0, 0, 0}
		m := Init[float64](10, 10)
		m = m.Slice(3, 2, 5, 3)
		for _, row := range m.IterRows() {
			if !reflect.DeepEqual(row, want) {
				fmt.Println(row)
				t.Errorf("iter wrong")
			}
		}
	})
	t.Run("testing sum", func(t *testing.T) {
		m := FromValues(2, 2, []float64{1, 1, 1, 1})
		got := m.Sum()
		if got != 4 {
			t.Errorf("Slice wrong")
		}
	})
	t.Run("testing mul", func(t *testing.T) {
		// need to test with different strides
		a := Matrix[float64]{Rows: 2, Cols: 2, Data: []float64{3, 3, 3, 3}, Stride: 2}
		b := Matrix[float64]{Rows: 2, Cols: 2, Data: []float64{4, 4, 4, 7, 7}, Stride: 3}
		c := a.Mul(b)
		for i, row := range c.IterRows() {
			if i == 0 {
				if !reflect.DeepEqual(row, Vec[float64]{12, 12}) {
					t.Errorf("mul wrong")
				}
			} else {
				if !reflect.DeepEqual(row, Vec[float64]{21, 21}) {
					t.Errorf("mul wrong")
				}
			}
		}
	})
}
