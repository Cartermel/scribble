package main

import "github.com/hajimehoshi/ebiten/v2"

// checks if the given key slices are equal (contain same keys).
// assumes the keys are in the same order
func KeysEqual(a, b []ebiten.Key) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

var (
	KeybindUndo = []ebiten.Key{ebiten.KeyZ, ebiten.KeyControlLeft, ebiten.KeyControl}
	KeybindDrag = []ebiten.Key{ebiten.KeySpace}
)
