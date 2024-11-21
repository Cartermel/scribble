package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// checks if the given key slices are equal (contain same keys).
// assumes the keys are in the same order
func keysEqual(a, b []ebiten.Key) bool {
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

// reads currently pressed keys, returning the specified keybind
func HandleKeyBindRead() Keybind {
	// TODO: use allocated slice for more efficient reads? dont rly care
	keys := inpututil.AppendPressedKeys(nil)
	if len(keys) == 0 {
		return KeybindNone
	}

	for k, v := range keyBindMap {
		if keysEqual(keys, v) {
			return k
		}
	}

	return KeybindNone
}

type Keybind int

var (
	KeybindNone Keybind = -1
	KeybindUndo Keybind = 0
	KeybindRedo Keybind = 1
	KeybindDrag Keybind = 2
)

// registry of all possible keybinds
var keyBindMap = map[Keybind][]ebiten.Key{
	// TODO: agnostic keys? im hard coding for left ctrl / shift because its exactly what ebiten returns
	KeybindUndo: []ebiten.Key{ebiten.KeyZ, ebiten.KeyControlLeft, ebiten.KeyControl},
	KeybindRedo: []ebiten.Key{ebiten.KeyZ, ebiten.KeyControlLeft, ebiten.KeyShiftLeft, ebiten.KeyControl, ebiten.KeyShift},
	KeybindDrag: []ebiten.Key{ebiten.KeySpace},
}
