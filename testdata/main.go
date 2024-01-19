package testdata

// @copyable
type YourStruct struct {
	// Struct fields definition
	Field1 string
	Field2 int
	field3 int
}

// @copyable
// @ptrrecv false
// @name clone
// exportedonly  true
type AnotherStruct struct {
	// Struct fields definition
	S  *YourStruct
	S2 int `deepcopy:"-"` // no export
}
