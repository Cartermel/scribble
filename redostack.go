package main

type RedoStack[T any] struct {
	items []T

	// always +1 than the actual index (ie, this == 1 when theres 1 item, so the real index is 0)
	// this tracks our position
	idx int
}

func (u *RedoStack[T]) Push(item T) {
	// if the idx is not the same as the length, it means some items have been popped
	// shrink the list to the current index and append.
	if u.idx != len(u.items) {
		u.items = u.items[0:u.idx] // TODO: plus 1?
	}

	u.items = append(u.items, item)
	u.idx++
}

// pops an item from the stack, returning it and a bool for whether or not an item was returned
func (u *RedoStack[T]) Pop() (T, bool) {
	u.idx--
	if u.idx < 0 {
		u.idx = 0
		return *new(T), false
	} else {
		return u.items[u.idx], true
	}
}

// Redo an action, re-applies the last `Pop` and returns it and a bool indicating if it could be redone
func (u *RedoStack[T]) Redo() (T, bool) {
	if u.idx < len(u.items) {
		val := u.items[u.idx]
		u.idx++
		return val, true
	}
	return *new(T), false
}
