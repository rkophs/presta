package system

import (
	"github.com/rkophs/presta/err"
)

type System interface {
	Push(a StackEntry) err.Error
	Pop() (StackEntry, err.Error)
	FetchS(offset int) (StackEntry, err.Error)
	FetchM(memAddr int) (StackEntry, err.Error)
	FetchR(id int) (StackEntry, err.Error)
	SetS(offset int, entry StackEntry) err.Error
	SetM(memAddr int, entry StackEntry) err.Error
	SetR(id int, entry StackEntry) err.Error
	New(memAddr int) err.Error
	Release(addr int) err.Error
	Shrink(result StackEntry) err.Error
	Goto(offset int) err.Error
	Expand() err.Error
	Exit() err.Error
}
