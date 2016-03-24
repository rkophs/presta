package system

type System interface {
	Push(a StackEntry)
	FetchS(offset int) StackEntry
	FetchM(memAddr int) StackEntry
	FetchR(id int) StackEntry
	SetS(offset int, entry StackEntry)
	SetM(memAddr int, entry StackEntry)
	SetR(id int, entry StackEntry)
	Release(addr int)
	Goto(offset int)
	Call(offset int)
	Return(result StackEntry)
	Exit(result StackEntry)
	SetError(e string)
}
