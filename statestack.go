package scribble

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// stack structure for keeping track of image states for undo / redo
type StateStack struct {
	// state is simply a slice of images, i was worried about the possible memory usage of this
	// but I spammed it in testing and it barely went up 2mb for 100 state layers, so should be fine for any modern pc
	// TODO: limit to like 1000 or some number?
	items []*ebiten.Image

	// always +1 than the actual index (ie, this == 1 when theres 1 item, so the real index is 0)
	// this tracks our position
	idx int
}

// pushes and item to the redo stack, only increments the index counter if passed `incrementIdx` is true
// leave `incrementIdx` false in order to push an un-poppable stack (ie, only redo-able)
func (u *StateStack) Push(item *ebiten.Image) {
	// if the idx is not the same as the length, it means some items have been popped
	// shrink the list to the current index and append.
	if u.idx != len(u.items) {
		u.items = u.items[0:u.idx] // TODO: plus 1?
	}

	u.items = append(u.items, item)
	u.idx++
}

// pops an item from the stack, returning it and a bool for whether or not an item was returned
// must pass the current state at time of undo-ing, for future re-dos
func (u *StateStack) Undo(currentState *ebiten.Image) (*ebiten.Image, bool) {
	// if we're undoing from the very top, push without incrementing the index
	if u.idx == len(u.items) {
		u.items = append(u.items, currentState)
	}

	u.idx--
	if u.idx < 0 {
		u.idx = 0
		return nil, false
	} else {
		return u.items[u.idx], true
	}
}

// Redo an action, re-applies the last `Pop` and returns it and a bool indicating if it could be redone
func (u *StateStack) Redo() (*ebiten.Image, bool) {
	if u.idx+1 < len(u.items) {
		u.idx++
		val := u.items[u.idx]
		return val, true
	}

	return nil, false
}
