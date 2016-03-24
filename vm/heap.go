package vm

import (
	"github.com/rkophs/presta/system"
)

type Heap struct {
	heap map[int]system.StackEntry
}

func NewHeap() *Heap {
	return &Heap{heap: make(map[int]system.StackEntry)}
}

func (h *Heap) Fetch(memAddr int) system.StackEntry {
	return h.heap[memAddr]
}

func (h *Heap) Set(memAddr int, entry system.StackEntry) {
	h.heap[memAddr] = entry
}

func (h *Heap) Release(addr int) {
	delete(h.heap, addr)
}
