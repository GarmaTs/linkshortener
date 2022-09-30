package lesson1

import (
	"testing"
)

func TestGreet(t *testing.T) {
	t.Run("greet Bob", func(t *testing.T) {
		got := Greet("Bob")
		want := "Hello Bob!"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
