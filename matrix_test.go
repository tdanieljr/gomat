package gomat

import (
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
	t.Run("testing get/set", func(t *testing.T) {
		m := Init[float64](20, 20)
		if m.Get(0, 0) != 0 {
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
}
